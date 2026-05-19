package rest

import (
	"github.com/eventhub/gateway/internal/client"
	"github.com/eventhub/gateway/pkg/auth"
)

// Handler serves REST endpoints documented via Swagger.
type Handler struct {
	Clients *client.GRPCClients
	JWT     *auth.JWTManager
}

// NewHandler creates a REST handler with gRPC clients and JWT manager.
func NewHandler(clients *client.GRPCClients, jwt *auth.JWTManager) *Handler {
	return &Handler{Clients: clients, JWT: jwt}
}
