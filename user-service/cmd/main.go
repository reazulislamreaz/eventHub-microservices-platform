package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/eventhub/pkg/database"
	"github.com/eventhub/pkg/logger"
	userv1 "github.com/eventhub/proto/gen/user/v1"
	"github.com/eventhub/user-service/config"
	"github.com/eventhub/user-service/internal/repository"
	"github.com/eventhub/user-service/internal/seed"
	userservice "github.com/eventhub/user-service/internal/service"
	usergrpc "github.com/eventhub/user-service/internal/transport/grpc"
	userdb "github.com/eventhub/user-service/pkg/database"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	_ = godotenv.Load()

	log, err := logger.New("user-service")
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("load config", zap.Error(err))
	}

	db, err := userdb.Connect(cfg.DB.DSN())
	if err != nil {
		log.Fatal("database connect", zap.Error(err))
	}
	if err := database.WaitForPostgres(db, log, 30); err != nil {
		log.Fatal("database ready", zap.Error(err))
	}
	if err := userdb.Migrate(db); err != nil {
		log.Fatal("database migrate", zap.Error(err))
	}

	repo := repository.NewUserRepository(db)
	svc := userservice.NewUserService(repo)
	seed.Admin(context.Background(), svc, log)
	handler := usergrpc.NewUserHandler(svc)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatal("listen", zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	userv1.RegisterUserServiceServer(grpcServer, handler)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("user.v1.UserService", healthpb.HealthCheckResponse_SERVING)
	reflection.Register(grpcServer)

	go func() {
		log.Info("user-service gRPC listening", zap.String("port", cfg.GRPCPort))
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("serve", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting down user-service")
	grpcServer.GracefulStop()
}
