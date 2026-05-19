package database

import (
	"fmt"

	"github.com/eventhub/ticket-service/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}
	return db, nil
}

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&model.Ticket{}); err != nil {
		return err
	}
	// Prevent duplicate confirmed bookings per user/event (race-safe with app check).
	return db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_tickets_user_event_confirmed
		ON tickets (user_id, event_id)
		WHERE status = 'confirmed'
	`).Error
}
