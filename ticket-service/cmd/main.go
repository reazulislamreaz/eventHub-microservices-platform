package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/eventhub/pkg/database"
	"github.com/eventhub/pkg/grpcutil"
	"github.com/eventhub/pkg/logger"
	eventv1 "github.com/eventhub/proto/gen/event/v1"
	ticketv1 "github.com/eventhub/proto/gen/ticket/v1"
	"github.com/eventhub/ticket-service/config"
	"github.com/eventhub/ticket-service/internal/repository"
	ticketservice "github.com/eventhub/ticket-service/internal/service"
	ticketgrpc "github.com/eventhub/ticket-service/internal/transport/grpc"
	ticketdb "github.com/eventhub/ticket-service/pkg/database"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	_ = godotenv.Load()

	log, err := logger.New("ticket-service")
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("load config", zap.Error(err))
	}

	db, err := ticketdb.Connect(cfg.DB.DSN())
	if err != nil {
		log.Fatal("database connect", zap.Error(err))
	}
	if err := database.WaitForPostgres(db, log, 30); err != nil {
		log.Fatal("database ready", zap.Error(err))
	}
	if err := ticketdb.Migrate(db); err != nil {
		log.Fatal("database migrate", zap.Error(err))
	}

	ctx := context.Background()
	if err := grpcutil.WaitForService(ctx, cfg.EventServiceAddr, "event.v1.EventService", log, 60); err != nil {
		log.Fatal("event-service not ready", zap.Error(err))
	}

	eventConn, err := grpc.NewClient(cfg.EventServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("event service client", zap.Error(err))
	}
	defer eventConn.Close()
	eventClient := eventv1.NewEventServiceClient(eventConn)

	repo := repository.NewTicketRepository(db)
	svc := ticketservice.NewTicketService(repo, eventClient)
	handler := ticketgrpc.NewTicketHandler(svc)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatal("listen", zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	ticketv1.RegisterTicketServiceServer(grpcServer, handler)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("ticket.v1.TicketService", healthpb.HealthCheckResponse_SERVING)
	reflection.Register(grpcServer)

	go func() {
		log.Info("ticket-service gRPC listening", zap.String("port", cfg.GRPCPort))
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("serve", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting down ticket-service")
	grpcServer.GracefulStop()
}
