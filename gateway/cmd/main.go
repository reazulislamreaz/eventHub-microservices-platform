package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/eventhub/gateway/config"
	"github.com/eventhub/gateway/internal/client"
	"github.com/eventhub/gateway/internal/graph"
	"github.com/eventhub/gateway/internal/handler/rest"
	"github.com/eventhub/gateway/internal/middleware"
	"github.com/eventhub/gateway/pkg/auth"
	"github.com/eventhub/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	_ "github.com/eventhub/gateway/docs"
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
// @host            localhost:8080
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

	router.Handle("/health", http.HandlerFunc(rest.Health)).Methods(http.MethodGet)
	router.Handle("/ready", http.HandlerFunc(rest.Ready)).Methods(http.MethodGet)
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	router.Handle("/", playground.Handler("EventHub GraphQL", "/query"))
	router.Handle("/query", srv)

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler(middleware.AuthMiddleware(jwtManager)(router))

	httpServer := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: handler,
	}

	go func() {
		log.Info("gateway listening",
			zap.String("http", cfg.HTTPPort),
			zap.String("graphql", "/query"),
			zap.String("playground", "/"),
			zap.String("swagger", "/swagger/index.html"),
			zap.String("api_docs", "/api/v1/docs"),
		)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("http serve", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting down gateway")
	_ = httpServer.Close()
}
