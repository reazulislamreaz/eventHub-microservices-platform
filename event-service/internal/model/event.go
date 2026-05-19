package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	StatusPublished = "published"
	StatusCancelled = "cancelled"
)

type Event struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Title          string    `gorm:"size:255;not null" json:"title"`
	Description    string    `gorm:"type:text" json:"description"`
	Location       string    `gorm:"size:255;not null" json:"location"`
	StartTime      time.Time `gorm:"not null" json:"start_time"`
	EndTime        time.Time `gorm:"not null" json:"end_time"`
	Capacity       int32     `gorm:"not null" json:"capacity"`
	AvailableSeats int32     `gorm:"not null" json:"available_seats"`
	Status         string    `gorm:"size:50;not null;default:published;index" json:"status"`
	CreatedBy      uuid.UUID `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (e *Event) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	if e.AvailableSeats == 0 && e.Capacity > 0 {
		e.AvailableSeats = e.Capacity
	}
	if e.Status == "" {
		e.Status = StatusPublished
	}
	return nil
}
