package graph

import (
	"github.com/eventhub/gateway/config"
	"github.com/eventhub/gateway/internal/client"
	"github.com/eventhub/gateway/pkg/auth"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies
// these resolvers require here.

type Resolver struct {
	Clients *client.GRPCClients
	JWT     *auth.JWTManager
	Config  *config.Config
}
