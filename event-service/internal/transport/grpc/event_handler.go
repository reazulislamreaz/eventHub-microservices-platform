package grpc

import (
	"context"
	"errors"
	"time"

	eventv1 "github.com/eventhub/proto/gen/event/v1"
	"github.com/eventhub/event-service/internal/model"
	"github.com/eventhub/event-service/internal/repository"
	"github.com/eventhub/event-service/internal/service"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EventHandler struct {
	eventv1.UnimplementedEventServiceServer
	svc service.EventService
}

func NewEventHandler(svc service.EventService) *EventHandler {
	return &EventHandler{svc: svc}
}

func (h *EventHandler) CreateEvent(ctx context.Context, req *eventv1.CreateEventRequest) (*eventv1.CreateEventResponse, error) {
	createdBy, err := uuid.Parse(req.GetCreatedBy())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid created_by")
	}
	start, err := time.Parse(time.RFC3339, req.GetStartTime())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid start_time")
	}
	end, err := time.Parse(time.RFC3339, req.GetEndTime())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid end_time")
	}

	event, err := h.svc.CreateEvent(ctx, req.GetTitle(), req.GetDescription(), req.GetLocation(), req.GetCategory(), req.GetPriceCents(), start, end, req.GetCapacity(), createdBy)
	if err != nil {
		return nil, mapError(err)
	}
	return &eventv1.CreateEventResponse{Event: toProtoEvent(event)}, nil
}

func (h *EventHandler) ListEvents(ctx context.Context, req *eventv1.ListEventsRequest) (*eventv1.ListEventsResponse, error) {
	out, err := h.svc.ListEvents(ctx, service.EventListInput{
		Page:     req.GetPage(),
		PageSize: req.GetPageSize(),
		Search:   req.GetSearch(),
		Location: req.GetLocation(),
		Status:   req.GetStatus(),
		Category: req.GetCategory(),
	})
	if err != nil {
		return nil, mapError(err)
	}
	resp := &eventv1.ListEventsResponse{
		Total:    out.Total,
		Page:     out.Page,
		PageSize: out.PageSize,
		Events:   make([]*eventv1.Event, 0, len(out.Events)),
	}
	for i := range out.Events {
		e := out.Events[i]
		resp.Events = append(resp.Events, toProtoEvent(&e))
	}
	return resp, nil
}

func (h *EventHandler) GetEvent(ctx context.Context, req *eventv1.GetEventRequest) (*eventv1.GetEventResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid event id")
	}
	event, err := h.svc.GetEvent(ctx, id)
	if err != nil {
		return nil, mapError(err)
	}
	return &eventv1.GetEventResponse{Event: toProtoEvent(event)}, nil
}

func (h *EventHandler) CancelEvent(ctx context.Context, req *eventv1.CancelEventRequest) (*eventv1.CancelEventResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid event id")
	}
	event, err := h.svc.CancelEvent(ctx, id)
	if err != nil {
		return nil, mapError(err)
	}
	return &eventv1.CancelEventResponse{Event: toProtoEvent(event)}, nil
}

func (h *EventHandler) ReserveSeat(ctx context.Context, req *eventv1.ReserveSeatRequest) (*eventv1.ReserveSeatResponse, error) {
	id, err := uuid.Parse(req.GetEventId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid event id")
	}
	event, err := h.svc.ReserveSeat(ctx, id)
	if err != nil {
		return nil, mapError(err)
	}
	return &eventv1.ReserveSeatResponse{Success: true, AvailableSeats: event.AvailableSeats}, nil
}

func (h *EventHandler) GetEventStats(ctx context.Context, _ *eventv1.GetEventStatsRequest) (*eventv1.GetEventStatsResponse, error) {
	st, err := h.svc.GetStats(ctx)
	if err != nil {
		return nil, mapError(err)
	}
	return &eventv1.GetEventStatsResponse{
		TotalEvents: st.TotalEvents, PublishedEvents: st.PublishedEvents,
		CancelledEvents: st.CancelledEvents, TotalCapacity: st.TotalCapacity,
		SeatsAvailable: st.SeatsAvailable,
	}, nil
}

func (h *EventHandler) ReleaseSeat(ctx context.Context, req *eventv1.ReleaseSeatRequest) (*eventv1.ReleaseSeatResponse, error) {
	id, err := uuid.Parse(req.GetEventId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid event id")
	}
	event, err := h.svc.ReleaseSeat(ctx, id)
	if err != nil {
		return nil, mapError(err)
	}
	return &eventv1.ReleaseSeatResponse{Success: true, AvailableSeats: event.AvailableSeats}, nil
}

func toProtoEvent(e *model.Event) *eventv1.Event {
	return &eventv1.Event{
		Id:             e.ID.String(),
		Title:          e.Title,
		Description:    e.Description,
		Location:       e.Location,
		StartTime:      e.StartTime.UTC().Format(time.RFC3339),
		EndTime:        e.EndTime.UTC().Format(time.RFC3339),
		Capacity:       e.Capacity,
		AvailableSeats: e.AvailableSeats,
		CreatedBy:      e.CreatedBy.String(),
		CreatedAt:      e.CreatedAt.UTC().Format(time.RFC3339),
		Status:         e.Status,
		Category:       e.Category,
		PriceCents:     e.PriceCents,
	}
}

func mapError(err error) error {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, repository.ErrEventNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, repository.ErrNoSeatsAvailable):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, repository.ErrEventCancelled):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, repository.ErrEventStarted):
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
