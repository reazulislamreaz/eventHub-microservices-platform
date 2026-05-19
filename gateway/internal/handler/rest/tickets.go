package rest

import (
	"encoding/json"
	"net/http"

	ticketv1 "github.com/eventhub/proto/gen/ticket/v1"
	"github.com/eventhub/gateway/pkg/auth"
	"github.com/gorilla/mux"
)

// BookTicket godoc
// @Summary      Book a ticket
// @Description  Books a ticket for the authenticated user. Reserves a seat on the event.
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      BookTicketRequest  true  "Event to book"
// @Success      201   {object}  Ticket
// @Failure      400   {object}  ErrorResponse
// @Failure      401   {object}  ErrorResponse
// @Failure      409   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /api/v1/tickets [post]
func (h *Handler) BookTicket(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	var req BookTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.EventID == "" {
		writeError(w, http.StatusBadRequest, "eventId is required")
		return
	}

	resp, err := h.Clients.Ticket.CreateTicket(r.Context(), &ticketv1.CreateTicketRequest{
		UserId:  claims.UserID,
		EventId: req.EventID,
	})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, mapProtoTicket(resp.GetTicket()))
}

// GetUserTickets godoc
// @Summary      Get user tickets
// @Description  Returns tickets for a user. Users may only access their own tickets unless admin.
// @Tags         tickets
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "User ID (UUID)"
// @Success      200  {array}   Ticket
// @Failure      401  {object}  ErrorResponse
// @Failure      403  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/users/{id}/tickets [get]
func (h *Handler) GetUserTickets(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	userID := mux.Vars(r)["id"]
	if claims.Role != "admin" && claims.UserID != userID {
		writeError(w, http.StatusForbidden, "cannot access another user's tickets")
		return
	}

	resp, err := h.Clients.Ticket.GetTicketsByUser(r.Context(), &ticketv1.GetTicketsByUserRequest{UserId: userID})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	tickets := make([]Ticket, 0, len(resp.GetTickets()))
	for _, t := range resp.GetTickets() {
		tickets = append(tickets, mapProtoTicket(t))
	}
	writeJSON(w, http.StatusOK, tickets)
}

// CancelTicket godoc
// @Summary      Cancel ticket
// @Description  Cancels a booking and releases the seat back to the event
// @Tags         tickets
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Ticket ID"
// @Success      200  {object}  Ticket
// @Router       /api/v1/tickets/{id}/cancel [post]
func (h *Handler) CancelTicket(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	ticketID := mux.Vars(r)["id"]
	resp, err := h.Clients.Ticket.CancelTicket(r.Context(), &ticketv1.CancelTicketRequest{
		Id: ticketID, UserId: claims.UserID,
	})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, mapProtoTicket(resp.GetTicket()))
}

func mapProtoTicket(t *ticketv1.Ticket) Ticket {
	ticket := Ticket{
		ID: t.GetId(), UserID: t.GetUserId(), EventID: t.GetEventId(),
		Status: t.GetStatus(), TicketCode: t.GetTicketCode(), CreatedAt: t.GetCreatedAt(),
	}
	if checked := t.GetCheckedInAt(); checked != "" {
		ticket.CheckedInAt = &checked
	}
	return ticket
}
