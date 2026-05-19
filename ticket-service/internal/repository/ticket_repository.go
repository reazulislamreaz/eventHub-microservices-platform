package repository

import (
	"context"
	"errors"

	"github.com/eventhub/ticket-service/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrTicketNotFound = errors.New("ticket not found")
var ErrDuplicateBooking = errors.New("user already has a ticket for this event")

type TicketRepository interface {
	Create(ctx context.Context, ticket *model.Ticket) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Ticket, error)
	GetByUser(ctx context.Context, userID uuid.UUID) ([]model.Ticket, error)
	ExistsForUserEvent(ctx context.Context, userID, eventID uuid.UUID) (bool, error)
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

func (r *ticketRepository) GetByUser(ctx context.Context, userID uuid.UUID) ([]model.Ticket, error) {
	var tickets []model.Ticket
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&tickets).Error
	return tickets, err
}

func (r *ticketRepository) ExistsForUserEvent(ctx context.Context, userID, eventID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Ticket{}).
		Where("user_id = ? AND event_id = ? AND status = ?", userID, eventID, model.StatusConfirmed).
		Count(&count).Error
	return count > 0, err
}
