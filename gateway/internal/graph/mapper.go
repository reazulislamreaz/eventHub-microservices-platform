package graph

import (
	eventv1 "github.com/eventhub/proto/gen/event/v1"
	ticketv1 "github.com/eventhub/proto/gen/ticket/v1"
	userv1 "github.com/eventhub/proto/gen/user/v1"
	"github.com/eventhub/gateway/internal/graph/model"
)

func mapUser(u *userv1.User) *model.User {
	if u == nil {
		return nil
	}
	return &model.User{
		ID:        u.GetId(),
		Email:     u.GetEmail(),
		Name:      u.GetName(),
		Role:      u.GetRole(),
		CreatedAt: u.GetCreatedAt(),
	}
}

func mapEvent(e *eventv1.Event) *model.Event {
	if e == nil {
		return nil
	}
	return &model.Event{
		ID:             e.GetId(),
		Title:          e.GetTitle(),
		Description:    e.GetDescription(),
		Location:       e.GetLocation(),
		StartTime:      e.GetStartTime(),
		EndTime:        e.GetEndTime(),
		Capacity:       int(e.GetCapacity()),
		AvailableSeats: int(e.GetAvailableSeats()),
		Status:         e.GetStatus(),
		CreatedBy:      e.GetCreatedBy(),
		CreatedAt:      e.GetCreatedAt(),
	}
}

func mapTicket(t *ticketv1.Ticket) *model.Ticket {
	if t == nil {
		return nil
	}
	return &model.Ticket{
		ID:         t.GetId(),
		UserID:     t.GetUserId(),
		EventID:    t.GetEventId(),
		Status:     t.GetStatus(),
		TicketCode: t.GetTicketCode(),
		CreatedAt:  t.GetCreatedAt(),
	}
}
