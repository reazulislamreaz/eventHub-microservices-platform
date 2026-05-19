package seed

import (
	"context"
	"errors"
	"os"

	"github.com/eventhub/user-service/internal/model"
	"github.com/eventhub/user-service/internal/service"
	"go.uber.org/zap"
)

// Admin seeds a default admin account when SEED_ADMIN=true.
func Admin(ctx context.Context, svc service.UserService, log *zap.Logger) {
	if os.Getenv("SEED_ADMIN") != "true" {
		return
	}
	_, err := svc.CreateUser(ctx, "admin@eventhub.io", "Platform Admin", "AdminPass123!", model.RoleAdmin)
	if err != nil {
		if errors.Is(err, service.ErrEmailTaken) {
			log.Info("admin account already exists", zap.String("email", "admin@eventhub.io"))
			return
		}
		log.Warn("admin seed failed", zap.Error(err))
		return
	}
	log.Info("seeded default admin", zap.String("email", "admin@eventhub.io"))
}
