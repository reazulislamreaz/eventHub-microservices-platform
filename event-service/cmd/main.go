package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eventhub/pkg/logger"
	eventv1 "github.com/eventhub/proto/gen/event/v1"
	"github.com/eventhub/event-service/config"
	"github.com/eventhub/event-service/internal/repository"
	eventservice "github.com/eventhub/event-service/internal/service"
	eventgrpc "github.com/eventhub/event-service/internal/transport/grpc"
	"github.com/eventhub/event-service/pkg/database"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()

	log, err := logger.New("event-service")
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("load config", zap.Error(err))
	}

	db, err := database.Connect(cfg.DB.DSN())
	if err != nil {
		log.Fatal("database connect", zap.Error(err))
	}
	if err := waitForDB(db, log); err != nil {
		log.Fatal("database ready", zap.Error(err))
	}
	if err := database.Migrate(db); err != nil {
		log.Fatal("database migrate", zap.Error(err))
	}

	repo := repository.NewEventRepository(db)
	svc := eventservice.NewEventService(repo)
	handler := eventgrpc.NewEventHandler(svc)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatal("listen", zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	eventv1.RegisterEventServiceServer(grpcServer, handler)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("event.v1.EventService", healthpb.HealthCheckResponse_SERVING)
	reflection.Register(grpcServer)

	go func() {
		log.Info("event-service gRPC listening", zap.String("port", cfg.GRPCPort))
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("serve", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting down event-service")
	grpcServer.GracefulStop()
}

func waitForDB(db *gorm.DB, log *zap.Logger) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	var lastErr error
	for i := 0; i < 30; i++ {
		if err := sqlDB.Ping(); err == nil {
			return nil
		} else {
			lastErr = err
		}
		log.Warn("waiting for database", zap.Int("attempt", i+1))
		time.Sleep(time.Second)
	}
	if lastErr != nil {
		return lastErr
	}
	return fmt.Errorf("database not ready")
}
