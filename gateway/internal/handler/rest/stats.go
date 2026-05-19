package rest

import (
	"net/http"

	eventv1 "github.com/eventhub/proto/gen/event/v1"
	ticketv1 "github.com/eventhub/proto/gen/ticket/v1"
	userv1 "github.com/eventhub/proto/gen/user/v1"
)

// PlatformStats aggregates dashboard metrics (admin).
type PlatformStats struct {
	Users   UserStats   `json:"users"`
	Events  EventStats  `json:"events"`
	Tickets TicketStats `json:"tickets"`
}

type UserStats struct {
	TotalUsers int32 `json:"totalUsers" example:"120"`
	AdminUsers int32 `json:"adminUsers" example:"2"`
}

type EventStats struct {
	TotalEvents     int32 `json:"totalEvents" example:"25"`
	PublishedEvents int32 `json:"publishedEvents" example:"20"`
	CancelledEvents int32 `json:"cancelledEvents" example:"5"`
	TotalCapacity   int32 `json:"totalCapacity" example:"5000"`
	SeatsAvailable  int32 `json:"seatsAvailable" example:"1200"`
}

type TicketStats struct {
	TotalTickets     int32 `json:"totalTickets" example:"800"`
	ConfirmedTickets int32 `json:"confirmedTickets" example:"650"`
	CheckedInTickets int32 `json:"checkedInTickets" example:"400"`
	WaitlistEntries  int32 `json:"waitlistEntries" example:"45"`
}

// GetPlatformStats godoc
// @Summary      Platform statistics
// @Description  Admin dashboard: users, events, and ticket metrics
// @Tags         admin
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  PlatformStats
// @Router       /api/v1/admin/stats [get]
func (h *Handler) GetPlatformStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	users, err := h.Clients.User.GetUserStats(ctx, &userv1.GetUserStatsRequest{})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	events, err := h.Clients.Event.GetEventStats(ctx, &eventv1.GetEventStatsRequest{})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	tickets, err := h.Clients.Ticket.GetTicketStats(ctx, &ticketv1.GetTicketStatsRequest{})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, PlatformStats{
		Users: UserStats{
			TotalUsers: users.GetTotalUsers(),
			AdminUsers: users.GetAdminUsers(),
		},
		Events: EventStats{
			TotalEvents: events.GetTotalEvents(), PublishedEvents: events.GetPublishedEvents(),
			CancelledEvents: events.GetCancelledEvents(), TotalCapacity: events.GetTotalCapacity(),
			SeatsAvailable: events.GetSeatsAvailable(),
		},
		Tickets: TicketStats{
			TotalTickets: tickets.GetTotalTickets(), ConfirmedTickets: tickets.GetConfirmedTickets(),
			CheckedInTickets: tickets.GetCheckedInTickets(), WaitlistEntries: tickets.GetWaitlistEntries(),
		},
	})
}
