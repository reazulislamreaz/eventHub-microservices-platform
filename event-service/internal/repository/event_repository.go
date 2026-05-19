package repository

import (
	"context"
	"errors"

	"github.com/eventhub/event-service/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrEventNotFound = errors.New("event not found")
var ErrNoSeatsAvailable = errors.New("no seats available")

type EventRepository interface {
	Create(ctx context.Context, event *model.Event) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Event, error)
	List(ctx context.Context) ([]model.Event, error)
	ReserveSeat(ctx context.Context, id uuid.UUID) (*model.Event, error)
	ReleaseSeat(ctx context.Context, id uuid.UUID) (*model.Event, error)
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) Create(ctx context.Context, event *model.Event) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *eventRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Event, error) {
	var event model.Event
	err := r.db.WithContext(ctx).First(&event, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEventNotFound
		}
		return nil, err
	}
	return &event, nil
}

func (r *eventRepository) List(ctx context.Context) ([]model.Event, error) {
	var events []model.Event
	err := r.db.WithContext(ctx).Order("start_time ASC").Find(&events).Error
	return events, err
}

func (r *eventRepository) ReserveSeat(ctx context.Context, id uuid.UUID) (*model.Event, error) {
	var event model.Event
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&event, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrEventNotFound
			}
			return err
		}
		if event.AvailableSeats <= 0 {
			return ErrNoSeatsAvailable
		}
		event.AvailableSeats--
		return tx.Save(&event).Error
	})
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventRepository) ReleaseSeat(ctx context.Context, id uuid.UUID) (*model.Event, error) {
	var event model.Event
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&event, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrEventNotFound
			}
			return err
		}
		if event.AvailableSeats < event.Capacity {
			event.AvailableSeats++
		}
		return tx.Save(&event).Error
	})
	if err != nil {
		return nil, err
	}
	return &event, nil
}
