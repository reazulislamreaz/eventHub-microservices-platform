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

// Common event categories for discovery pages.
const (
	CategoryMusic       = "music"
	CategoryTech        = "tech"
	CategorySports      = "sports"
	CategoryConference  = "conference"
	CategoryWorkshop    = "workshop"
	CategoryOther       = "other"
)

type Event struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Title          string    `gorm:"size:255;not null" json:"title"`
	Description    string    `gorm:"type:text" json:"description"`
	Location       string    `gorm:"size:255;not null" json:"location"`
	Category       string    `gorm:"size:50;not null;default:other;index" json:"category"`
	PriceCents     int64     `gorm:"not null;default:0" json:"price_cents"`
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
	if e.Category == "" {
		e.Category = CategoryOther
	}
	return nil
}
