package grpc

import (
	"context"
	"errors"

	ticketv1 "github.com/eventhub/proto/gen/ticket/v1"
	"github.com/eventhub/ticket-service/internal/model"
	"github.com/eventhub/ticket-service/internal/repository"
	"github.com/eventhub/ticket-service/internal/service"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TicketHandler struct {
	ticketv1.UnimplementedTicketServiceServer
	svc service.TicketService
}

func NewTicketHandler(svc service.TicketService) *TicketHandler {
	return &TicketHandler{svc: svc}
}

func (h *TicketHandler) CreateTicket(ctx context.Context, req *ticketv1.CreateTicketRequest) (*ticketv1.CreateTicketResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}
	eventID, err := uuid.Parse(req.GetEventId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid event id")
	}

	ticket, err := h.svc.CreateTicket(ctx, userID, eventID)
	if err != nil {
		return nil, mapError(err)
	}
	return &ticketv1.CreateTicketResponse{Ticket: toProtoTicket(ticket)}, nil
}

func (h *TicketHandler) CancelTicket(ctx context.Context, req *ticketv1.CancelTicketRequest) (*ticketv1.CancelTicketResponse, error) {
	ticketID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid ticket id")
	}
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}
	ticket, err := h.svc.CancelTicket(ctx, ticketID, userID)
	if err != nil {
		return nil, mapError(err)
	}
	return &ticketv1.CancelTicketResponse{Ticket: toProtoTicket(ticket)}, nil
}

func (h *TicketHandler) GetTicketsByUser(ctx context.Context, req *ticketv1.GetTicketsByUserRequest) (*ticketv1.GetTicketsByUserResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}
	tickets, err := h.svc.GetTicketsByUser(ctx, userID)
	if err != nil {
		return nil, mapError(err)
	}
	resp := &ticketv1.GetTicketsByUserResponse{Tickets: make([]*ticketv1.Ticket, 0, len(tickets))}
	for i := range tickets {
		t := tickets[i]
		resp.Tickets = append(resp.Tickets, toProtoTicket(&t))
	}
	return resp, nil
}

func (h *TicketHandler) GetTicket(ctx context.Context, req *ticketv1.GetTicketRequest) (*ticketv1.GetTicketResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid ticket id")
	}
	ticket, err := h.svc.GetTicket(ctx, id)
	if err != nil {
		return nil, mapError(err)
	}
	return &ticketv1.GetTicketResponse{Ticket: toProtoTicket(ticket)}, nil
}

func toProtoTicket(t *model.Ticket) *ticketv1.Ticket {
	return &ticketv1.Ticket{
		Id:         t.ID.String(),
		UserId:     t.UserID.String(),
		EventId:    t.EventID.String(),
		Status:     t.Status,
		TicketCode: t.TicketCode,
		CreatedAt:  t.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
	}
}

func mapError(err error) error {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, service.ErrDuplicateBooking):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, service.ErrSeatNotAvailable):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, repository.ErrTicketNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, repository.ErrTicketNotCancellable):
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
