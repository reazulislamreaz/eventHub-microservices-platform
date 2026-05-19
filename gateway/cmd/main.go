package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/eventhub/gateway/config"
	"github.com/eventhub/gateway/internal/client"
	"github.com/eventhub/pkg/grpcutil"
	"github.com/eventhub/gateway/internal/graph"
	"github.com/eventhub/gateway/internal/handler/rest"
	"github.com/eventhub/gateway/internal/middleware"
	gwmiddleware "github.com/eventhub/gateway/pkg/middleware"
	"github.com/eventhub/gateway/pkg/auth"
	"github.com/eventhub/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

)

// @title           EventHub Gateway API
// @version         1.0
// @description     EventHub microservices platform API gateway.
// @description     **REST API** — documented below; use Swagger UI at `/swagger/index.html` or `/docs`.
// @description     **GraphQL** — Playground at `/`, endpoint `POST /query`, schema at `/api/v1/graphql/schema`.
// @description
// @description     ## Authentication
// @description     Obtain a JWT via `POST /api/v1/auth/login` or `POST /api/v1/auth/register`, then pass `Authorization: Bearer <token>` for protected routes.
// @description
// @description     ## Default admin (Docker)
// @description     Email: admin@eventhub.io | Password: AdminPass123!
// @termsOfService  https://github.com/eventhub/platform
// @contact.name    EventHub API Support
// @contact.email   support@eventhub.io
// @license.name    MIT
// @license.url     https://opensource.org/licenses/MIT
// @BasePath        /
// @schemes         http https
// @tag.name        auth
// @tag.description User registration and login
// @tag.name        users
// @tag.description User profiles
// @tag.name        events
// @tag.description Event management
// @tag.name        tickets
// @tag.description Ticket booking
// @tag.name        health
// @tag.description Health and readiness probes
// @tag.name        documentation
// @tag.description API documentation endpoints
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description     JWT token. Format: Bearer {token}
func main() {
	_ = godotenv.Load()

	log, err := logger.New("gateway")
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	cfg := config.Load()
	configureSwaggerHost(cfg.SwaggerHost)

	startupCtx, startupCancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer startupCancel()
	for _, dep := range []struct {
		addr, name string
	}{
		{cfg.UserServiceAddr, "user.v1.UserService"},
		{cfg.EventServiceAddr, "event.v1.EventService"},
		{cfg.TicketServiceAddr, "ticket.v1.TicketService"},
	} {
		if err := grpcutil.WaitForService(startupCtx, dep.addr, dep.name, log, 60); err != nil {
			log.Fatal("dependency not ready", zap.String("service", dep.name), zap.Error(err))
		}
	}

	clients, err := client.NewGRPCClients(cfg.UserServiceAddr, cfg.EventServiceAddr, cfg.TicketServiceAddr)
	if err != nil {
		log.Fatal("grpc clients", zap.Error(err))
	}
	defer clients.Close()

	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiration)
	resolver := &graph.Resolver{Clients: clients, JWT: jwtManager, Config: cfg}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: resolver,
		Directives: graph.DirectiveRoot{
			Auth:    middleware.DirectiveAuth,
			HasRole: middleware.DirectiveHasRole,
		},
	}))

	router := mux.NewRouter()
	restHandler := rest.NewHandler(clients, jwtManager)
	rest.RegisterRoutes(router, restHandler)

	router.HandleFunc("/metrics", rest.Metrics).Methods(http.MethodGet)
	router.Handle("/health", http.HandlerFunc(rest.Health)).Methods(http.MethodGet)
	router.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		rest.ReadyWithDeps(w, r, map[string]string{
			"user":   cfg.UserServiceAddr,
			"event":  cfg.EventServiceAddr,
			"ticket": cfg.TicketServiceAddr,
		}, log)
	}).Methods(http.MethodGet)
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
	router.Handle("/", playground.Handler("EventHub GraphQL", "/query"))
	router.Handle("/query", srv)

	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders: []string{
			"Authorization",
			"Content-Type",
			"Accept",
			"Origin",
			"X-Requested-With",
			"X-CSRF-Token",
		},
		ExposedHeaders:   []string{"Content-Length", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}).Handler(
		gwmiddleware.RequestID(
			gwmiddleware.Metrics(
				gwmiddleware.RateLimit(20, 40)(
					middleware.AuthMiddleware(jwtManager)(router),
				),
			),
		),
	)

	httpServer := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Info("gateway listening",
			zap.String("http", cfg.HTTPPort),
			zap.String("graphql", "/query"),
			zap.String("playground", "/"),
			zap.String("api_docs", "/api/docs"),
			zap.String("swagger", "/swagger/index.html"),
		)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("http serve", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting down gateway")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error("shutdown", zap.Error(err))
	}
}
