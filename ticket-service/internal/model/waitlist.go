package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WaitlistEntry struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	EventID   uuid.UUID `gorm:"type:uuid;not null;index" json:"event_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (w *WaitlistEntry) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}
