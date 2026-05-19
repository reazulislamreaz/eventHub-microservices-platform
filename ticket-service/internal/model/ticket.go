package model

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	StatusConfirmed  = "confirmed"
	StatusCancelled  = "cancelled"
	StatusCheckedIn  = "checked_in"
)

type Ticket struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	EventID    uuid.UUID `gorm:"type:uuid;not null;index" json:"event_id"`
	Status     string    `gorm:"size:50;not null;default:confirmed" json:"status"`
	TicketCode  string     `gorm:"uniqueIndex;size:32;not null" json:"ticket_code"`
	CheckedInAt *time.Time `json:"checked_in_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (t *Ticket) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	if t.TicketCode == "" {
		code, err := generateTicketCode()
		if err != nil {
			return err
		}
		t.TicketCode = code
	}
	if t.Status == "" {
		t.Status = StatusConfirmed
	}
	return nil
}

func generateTicketCode() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return "EH-" + hex.EncodeToString(b), nil
}
