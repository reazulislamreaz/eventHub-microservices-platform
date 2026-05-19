package service

import (
	"context"
	"errors"
	"strings"

	eventv1 "github.com/eventhub/proto/gen/event/v1"
	"github.com/eventhub/ticket-service/internal/model"
	"github.com/eventhub/ticket-service/internal/repository"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidInput     = errors.New("invalid input")
	ErrDuplicateBooking = errors.New("user already booked this event")
	ErrSeatNotAvailable = errors.New("no seats available")
	ErrForbidden        = errors.New("forbidden")
	ErrAlreadyOnWaitlist = errors.New("already on waitlist for this event")
)

type TicketStatsOutput struct {
	TotalTickets     int32
	ConfirmedTickets int32
	CheckedInTickets int32
	WaitlistEntries  int32
}

type TicketService interface {
	CreateTicket(ctx context.Context, userID, eventID uuid.UUID) (*model.Ticket, error)
	CancelTicket(ctx context.Context, ticketID, userID uuid.UUID) (*model.Ticket, error)
	GetTicketsByUser(ctx context.Context, userID uuid.UUID) ([]model.Ticket, error)
	GetTicket(ctx context.Context, id uuid.UUID) (*model.Ticket, error)
	GetTicketByCode(ctx context.Context, code string) (*model.Ticket, error)
	CheckInTicket(ctx context.Context, code string) (*model.Ticket, error)
	JoinWaitlist(ctx context.Context, userID, eventID uuid.UUID) (*model.WaitlistEntry, error)
	GetStats(ctx context.Context) (*TicketStatsOutput, error)
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
		if isDuplicateKey(err) {
			_, _ = s.eventClient.ReleaseSeat(ctx, &eventv1.ReleaseSeatRequest{EventId: eventID.String()})
			return nil, ErrDuplicateBooking
		}
		if _, releaseErr := s.eventClient.ReleaseSeat(ctx, &eventv1.ReleaseSeatRequest{EventId: eventID.String()}); releaseErr != nil {
			return nil, errors.Join(err, releaseErr)
		}
		return nil, err
	}
	return ticket, nil
}

func (s *ticketService) CancelTicket(ctx context.Context, ticketID, userID uuid.UUID) (*model.Ticket, error) {
	if ticketID == uuid.Nil || userID == uuid.Nil {
		return nil, ErrInvalidInput
	}
	ticket, err := s.repo.Cancel(ctx, ticketID, userID)
	if err != nil {
		return nil, err
	}
	_, _ = s.eventClient.ReleaseSeat(ctx, &eventv1.ReleaseSeatRequest{EventId: ticket.EventID.String()})
	return ticket, nil
}

func (s *ticketService) GetTicketsByUser(ctx context.Context, userID uuid.UUID) ([]model.Ticket, error) {
	return s.repo.GetByUser(ctx, userID)
}

func (s *ticketService) GetTicket(ctx context.Context, id uuid.UUID) (*model.Ticket, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ticketService) GetTicketByCode(ctx context.Context, code string) (*model.Ticket, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return nil, ErrInvalidInput
	}
	return s.repo.GetByCode(ctx, code)
}

func (s *ticketService) CheckInTicket(ctx context.Context, code string) (*model.Ticket, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return nil, ErrInvalidInput
	}
	return s.repo.CheckIn(ctx, code)
}

func (s *ticketService) JoinWaitlist(ctx context.Context, userID, eventID uuid.UUID) (*model.WaitlistEntry, error) {
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

	onWaitlist, err := s.repo.WaitlistExists(ctx, userID, eventID)
	if err != nil {
		return nil, err
	}
	if onWaitlist {
		return nil, ErrAlreadyOnWaitlist
	}

	// Verify event exists and is bookable context (not cancelled).
	if _, err := s.eventClient.GetEvent(ctx, &eventv1.GetEventRequest{Id: eventID.String()}); err != nil {
		return nil, err
	}

	entry := &model.WaitlistEntry{UserID: userID, EventID: eventID}
	if err := s.repo.JoinWaitlist(ctx, entry); err != nil {
		if isDuplicateKey(err) {
			return nil, ErrAlreadyOnWaitlist
		}
		return nil, err
	}
	return entry, nil
}

func (s *ticketService) GetStats(ctx context.Context) (*TicketStatsOutput, error) {
	st, err := s.repo.Stats(ctx)
	if err != nil {
		return nil, err
	}
	return &TicketStatsOutput{
		TotalTickets: int32(st.TotalTickets), ConfirmedTickets: int32(st.ConfirmedTickets),
		CheckedInTickets: int32(st.CheckedInTickets), WaitlistEntries: int32(st.WaitlistEntries),
	}, nil
}

func isDuplicateKey(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate key") || strings.Contains(msg, "unique constraint")
}
