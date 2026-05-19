package service

import (
	"context"
	"errors"

	eventv1 "github.com/eventhub/proto/gen/event/v1"
	"github.com/eventhub/ticket-service/internal/model"
	"github.com/eventhub/ticket-service/internal/repository"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrDuplicateBooking   = errors.New("user already booked this event")
	ErrSeatNotAvailable   = errors.New("no seats available")
)

type TicketService interface {
	CreateTicket(ctx context.Context, userID, eventID uuid.UUID) (*model.Ticket, error)
	GetTicketsByUser(ctx context.Context, userID uuid.UUID) ([]model.Ticket, error)
	GetTicket(ctx context.Context, id uuid.UUID) (*model.Ticket, error)
}

type ticketService struct {
	repo        repository.TicketRepository
	eventClient eventv1.EventServiceClient
}

func NewTicketService(repo repository.TicketRepository, eventClient eventv1.EventServiceClient) TicketService {
	return &ticketService{repo: repo, eventClient: eventClient}
}

func (s *ticketService) CreateTicket(ctx context.Context, userID, eventID uuid.UUID) (*model.Ticket, error) {
	if userID == uuid.Nil || eventID == uuid.Nil {
		return nil, ErrInvalidInput
	}

	exists, err := s.repo.ExistsForUserEvent(ctx, userID, eventID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrDuplicateBooking
	}

	reserveResp, err := s.eventClient.ReserveSeat(ctx, &eventv1.ReserveSeatRequest{EventId: eventID.String()})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.FailedPrecondition {
				return nil, ErrSeatNotAvailable
			}
		}
		return nil, err
	}
	if !reserveResp.GetSuccess() {
		return nil, ErrSeatNotAvailable
	}

	ticket := &model.Ticket{
		UserID:  userID,
		EventID: eventID,
		Status:  model.StatusConfirmed,
	}
	if err := s.repo.Create(ctx, ticket); err != nil {
		_, _ = s.eventClient.ReleaseSeat(ctx, &eventv1.ReleaseSeatRequest{EventId: eventID.String()})
		return nil, err
	}
	return ticket, nil
}

func (s *ticketService) GetTicketsByUser(ctx context.Context, userID uuid.UUID) ([]model.Ticket, error) {
	return s.repo.GetByUser(ctx, userID)
}

func (s *ticketService) GetTicket(ctx context.Context, id uuid.UUID) (*model.Ticket, error) {
	return s.repo.GetByID(ctx, id)
}
