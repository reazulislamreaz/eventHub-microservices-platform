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
// @description     GraphQL gateway and REST utilities for EventHub microservices platform.
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter "Bearer {token}"
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
