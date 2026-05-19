package database

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// WaitForPostgres pings the database until ready or attempts are exhausted.
func WaitForPostgres(db *gorm.DB, log *zap.Logger, attempts int) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	var lastErr error
	for i := 0; i < attempts; i++ {
		if err := sqlDB.Ping(); err == nil {
			return nil
		}
		lastErr = err
		if log != nil {
			log.Warn("waiting for database", zap.Int("attempt", i+1))
		}
		time.Sleep(time.Second)
	}
	if lastErr != nil {
		return lastErr
	}
	return fmt.Errorf("database not ready after %d attempts", attempts)
}
