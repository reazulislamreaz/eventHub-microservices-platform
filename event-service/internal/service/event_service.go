package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/eventhub/event-service/internal/model"
	"github.com/eventhub/event-service/internal/repository"
	"github.com/google/uuid"
)

var ErrInvalidInput = errors.New("invalid input")

type EventService interface {
	CreateEvent(ctx context.Context, title, description, location string, startTime, endTime time.Time, capacity int32, createdBy uuid.UUID) (*model.Event, error)
	ListEvents(ctx context.Context) ([]model.Event, error)
	GetEvent(ctx context.Context, id uuid.UUID) (*model.Event, error)
	ReserveSeat(ctx context.Context, id uuid.UUID) (*model.Event, error)
	ReleaseSeat(ctx context.Context, id uuid.UUID) (*model.Event, error)
}

type eventService struct {
	repo repository.EventRepository
}

func NewEventService(repo repository.EventRepository) EventService {
	return &eventService{repo: repo}
}

func (s *eventService) CreateEvent(ctx context.Context, title, description, location string, startTime, endTime time.Time, capacity int32, createdBy uuid.UUID) (*model.Event, error) {
	title = strings.TrimSpace(title)
	location = strings.TrimSpace(location)
	if title == "" || location == "" || capacity <= 0 || !endTime.After(startTime) {
		return nil, ErrInvalidInput
	}

	event := &model.Event{
		Title:          title,
		Description:    description,
		Location:       location,
		StartTime:      startTime,
		EndTime:        endTime,
		Capacity:       capacity,
		AvailableSeats: capacity,
		CreatedBy:      createdBy,
	}
	if err := s.repo.Create(ctx, event); err != nil {
		return nil, err
	}
	return event, nil
}

func (s *eventService) ListEvents(ctx context.Context) ([]model.Event, error) {
	return s.repo.List(ctx)
}

func (s *eventService) GetEvent(ctx context.Context, id uuid.UUID) (*model.Event, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *eventService) ReserveSeat(ctx context.Context, id uuid.UUID) (*model.Event, error) {
	return s.repo.ReserveSeat(ctx, id)
}

func (s *eventService) ReleaseSeat(ctx context.Context, id uuid.UUID) (*model.Event, error) {
	return s.repo.ReleaseSeat(ctx, id)
}
