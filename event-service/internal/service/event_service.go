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

type EventListInput struct {
	Page     int32
	PageSize int32
	Search   string
	Location string
	Status   string
}

type EventListOutput struct {
	Events   []model.Event
	Total    int32
	Page     int32
	PageSize int32
}

type EventService interface {
	CreateEvent(ctx context.Context, title, description, location string, startTime, endTime time.Time, capacity int32, createdBy uuid.UUID) (*model.Event, error)
	ListEvents(ctx context.Context, in EventListInput) (*EventListOutput, error)
	GetEvent(ctx context.Context, id uuid.UUID) (*model.Event, error)
	CancelEvent(ctx context.Context, id uuid.UUID) (*model.Event, error)
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
		Status:         model.StatusPublished,
		CreatedBy:      createdBy,
	}
	if err := s.repo.Create(ctx, event); err != nil {
		return nil, err
	}
	return event, nil
}

func (s *eventService) ListEvents(ctx context.Context, in EventListInput) (*EventListOutput, error) {
	result, err := s.repo.List(ctx, repository.EventFilter{
		Page:     int(in.Page),
		PageSize: int(in.PageSize),
		Search:   in.Search,
		Location: in.Location,
		Status:   in.Status,
	})
	if err != nil {
		return nil, err
	}
	page := in.Page
	if page < 1 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	return &EventListOutput{
		Events:   result.Events,
		Total:    int32(result.Total),
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *eventService) GetEvent(ctx context.Context, id uuid.UUID) (*model.Event, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *eventService) CancelEvent(ctx context.Context, id uuid.UUID) (*model.Event, error) {
	return s.repo.Cancel(ctx, id)
}

func (s *eventService) ReserveSeat(ctx context.Context, id uuid.UUID) (*model.Event, error) {
	return s.repo.ReserveSeat(ctx, id)
}

func (s *eventService) ReleaseSeat(ctx context.Context, id uuid.UUID) (*model.Event, error) {
	return s.repo.ReleaseSeat(ctx, id)
}
