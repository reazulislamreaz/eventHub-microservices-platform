package repository

import (
	"context"
	"errors"
	"time"

	"github.com/eventhub/ticket-service/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrTicketNotFound       = errors.New("ticket not found")
	ErrDuplicateBooking     = errors.New("user already has a ticket for this event")
	ErrTicketNotCancellable = errors.New("ticket cannot be cancelled")
	ErrAlreadyCheckedIn     = errors.New("ticket already checked in")
	ErrInvalidTicketStatus  = errors.New("ticket cannot be checked in")
	ErrWaitlistExists       = errors.New("already on waitlist for this event")
	ErrWaitlistNotFound     = errors.New("waitlist entry not found")
)

type TicketStats struct {
	TotalTickets     int64
	ConfirmedTickets int64
	CheckedInTickets int64
	WaitlistEntries  int64
}

type TicketRepository interface {
	Create(ctx context.Context, ticket *model.Ticket) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Ticket, error)
	GetByCode(ctx context.Context, code string) (*model.Ticket, error)
	GetByUser(ctx context.Context, userID uuid.UUID) ([]model.Ticket, error)
	ExistsForUserEvent(ctx context.Context, userID, eventID uuid.UUID) (bool, error)
	Cancel(ctx context.Context, id, userID uuid.UUID) (*model.Ticket, error)
	CheckIn(ctx context.Context, code string) (*model.Ticket, error)
	Stats(ctx context.Context) (*TicketStats, error)
	JoinWaitlist(ctx context.Context, entry *model.WaitlistEntry) error
	WaitlistExists(ctx context.Context, userID, eventID uuid.UUID) (bool, error)
}

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) Create(ctx context.Context, ticket *model.Ticket) error {
	return r.db.WithContext(ctx).Create(ticket).Error
}

func (r *ticketRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Ticket, error) {
	var ticket model.Ticket
	err := r.db.WithContext(ctx).First(&ticket, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrTicketNotFound
	}
	return &ticket, err
}

func (r *ticketRepository) GetByCode(ctx context.Context, code string) (*model.Ticket, error) {
	var ticket model.Ticket
	err := r.db.WithContext(ctx).First(&ticket, "ticket_code = ?", code).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrTicketNotFound
	}
	return &ticket, err
}

func (r *ticketRepository) GetByUser(ctx context.Context, userID uuid.UUID) ([]model.Ticket, error) {
	var tickets []model.Ticket
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&tickets).Error
	return tickets, err
}

func (r *ticketRepository) ExistsForUserEvent(ctx context.Context, userID, eventID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Ticket{}).
		Where("user_id = ? AND event_id = ? AND status IN ?", userID, eventID, []string{model.StatusConfirmed, model.StatusCheckedIn}).
		Count(&count).Error
	return count > 0, err
}

func (r *ticketRepository) Cancel(ctx context.Context, id, userID uuid.UUID) (*model.Ticket, error) {
	var ticket model.Ticket
	err := r.db.WithContext(ctx).First(&ticket, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrTicketNotFound
	}
	if err != nil {
		return nil, err
	}
	if ticket.UserID != userID {
		return nil, ErrTicketNotFound
	}
	if ticket.Status != model.StatusConfirmed {
		return nil, ErrTicketNotCancellable
	}
	ticket.Status = model.StatusCancelled
	if err := r.db.WithContext(ctx).Save(&ticket).Error; err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *ticketRepository) CheckIn(ctx context.Context, code string) (*model.Ticket, error) {
	var ticket model.Ticket
	err := r.db.WithContext(ctx).First(&ticket, "ticket_code = ?", code).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrTicketNotFound
	}
	if err != nil {
		return nil, err
	}
	if ticket.Status == model.StatusCheckedIn {
		return nil, ErrAlreadyCheckedIn
	}
	if ticket.Status != model.StatusConfirmed {
		return nil, ErrInvalidTicketStatus
	}
	now := time.Now().UTC()
	ticket.Status = model.StatusCheckedIn
	ticket.CheckedInAt = &now
	if err := r.db.WithContext(ctx).Save(&ticket).Error; err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *ticketRepository) Stats(ctx context.Context) (*TicketStats, error) {
	var stats TicketStats
	if err := r.db.WithContext(ctx).Model(&model.Ticket{}).Count(&stats.TotalTickets).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Model(&model.Ticket{}).Where("status = ?", model.StatusConfirmed).Count(&stats.ConfirmedTickets).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Model(&model.Ticket{}).Where("status = ?", model.StatusCheckedIn).Count(&stats.CheckedInTickets).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Model(&model.WaitlistEntry{}).Count(&stats.WaitlistEntries).Error; err != nil {
		return nil, err
	}
	return &stats, nil
}

func (r *ticketRepository) JoinWaitlist(ctx context.Context, entry *model.WaitlistEntry) error {
	return r.db.WithContext(ctx).Create(entry).Error
}

func (r *ticketRepository) WaitlistExists(ctx context.Context, userID, eventID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.WaitlistEntry{}).
		Where("user_id = ? AND event_id = ?", userID, eventID).
		Count(&count).Error
	return count > 0, err
}
