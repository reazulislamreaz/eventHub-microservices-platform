package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/eventhub/event-service/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrEventNotFound = errors.New("event not found")
var ErrNoSeatsAvailable = errors.New("no seats available")
var ErrEventCancelled = errors.New("event is cancelled")

type EventFilter struct {
	Page     int
	PageSize int
	Search   string
	Location string
	Status   string
}

type EventListResult struct {
	Events []model.Event
	Total  int64
}

type EventRepository interface {
	Create(ctx context.Context, event *model.Event) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Event, error)
	List(ctx context.Context, filter EventFilter) (*EventListResult, error)
	Cancel(ctx context.Context, id uuid.UUID) (*model.Event, error)
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

func (r *eventRepository) List(ctx context.Context, filter EventFilter) (*EventListResult, error) {
	q := r.db.WithContext(ctx).Model(&model.Event{})
	if filter.Status != "" {
		q = q.Where("status = ?", filter.Status)
	} else {
		q = q.Where("status = ?", model.StatusPublished)
	}
	if filter.Location != "" {
		q = q.Where("location ILIKE ?", "%"+filter.Location+"%")
	}
	if filter.Search != "" {
		term := "%" + strings.ToLower(filter.Search) + "%"
		q = q.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ?", term, term)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	var events []model.Event
	err := q.Order("start_time ASC").Offset(offset).Limit(pageSize).Find(&events).Error
	return &EventListResult{Events: events, Total: total}, err
}

func (r *eventRepository) Cancel(ctx context.Context, id uuid.UUID) (*model.Event, error) {
	var event model.Event
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&event, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrEventNotFound
			}
			return err
		}
		if event.Status == model.StatusCancelled {
			return nil
		}
		event.Status = model.StatusCancelled
		return tx.Save(&event).Error
	})
	if err != nil {
		return nil, err
	}
	return &event, nil
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
		if event.Status == model.StatusCancelled {
			return ErrEventCancelled
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
