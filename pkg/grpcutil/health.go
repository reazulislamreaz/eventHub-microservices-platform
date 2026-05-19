package grpcutil

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// WaitForService dials addr and waits until the gRPC health check reports SERVING.
func WaitForService(ctx context.Context, addr, serviceName string, log *zap.Logger, attempts int) error {
	var lastErr error
	for i := 0; i < attempts; i++ {
		if err := checkOnce(ctx, addr, serviceName); err == nil {
			return nil
		} else {
			lastErr = err
		}
		if log != nil {
			log.Warn("waiting for gRPC service",
				zap.String("addr", addr),
				zap.String("service", serviceName),
				zap.Int("attempt", i+1),
			)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
		}
	}
	if lastErr != nil {
		return fmt.Errorf("service %s at %s not ready: %w", serviceName, addr, lastErr)
	}
	return fmt.Errorf("service %s at %s not ready", serviceName, addr)
}

func checkOnce(ctx context.Context, addr, serviceName string) error {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := healthpb.NewHealthClient(conn)
	checkCtx, checkCancel := context.WithTimeout(ctx, 3*time.Second)
	defer checkCancel()

	resp, err := client.Check(checkCtx, &healthpb.HealthCheckRequest{Service: serviceName})
	if err != nil {
		return err
	}
	if resp.GetStatus() != healthpb.HealthCheckResponse_SERVING {
		return fmt.Errorf("status %s", resp.GetStatus().String())
	}
	return nil
}
