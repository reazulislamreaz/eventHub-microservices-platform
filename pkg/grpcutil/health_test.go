package grpcutil

import (
	"context"
	"testing"
	"time"
)

func TestWaitForService_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := WaitForService(ctx, "127.0.0.1:59999", "test.Service", nil, 5)
	if err == nil {
		t.Fatal("expected error for unreachable service")
	}
}
