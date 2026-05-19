package rest

import (
	"encoding/json"
	"net/http"

	ticketv1 "github.com/eventhub/proto/gen/ticket/v1"
	"github.com/eventhub/gateway/pkg/auth"
	"github.com/gorilla/mux"
)

// CheckInRequest is the body for ticket check-in.
type CheckInRequest struct {
	TicketCode string `json:"ticketCode" example:"EH-a1b2c3d4e5f67890"`
}

// JoinWaitlistRequest joins the waitlist when sold out.
type JoinWaitlistRequest struct {
	EventID string `json:"eventId" example:"550e8400-e29b-41d4-a716-446655440001"`
}

// WaitlistEntry represents a waitlist row.
type WaitlistEntry struct {
	ID        string `json:"id"`
	UserID    string `json:"userId"`
	EventID   string `json:"eventId"`
	CreatedAt string `json:"createdAt"`
}

// VerifyTicket godoc
// @Summary      Verify ticket by code
// @Description  Look up a ticket by its QR/barcode code (owner or admin)
// @Tags         tickets
// @Produce      json
// @Security     BearerAuth
// @Param        code  query  string  true  "Ticket code"
// @Success      200   {object}  Ticket
// @Router       /api/v1/tickets/verify [get]
func (h *Handler) VerifyTicket(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	code := r.URL.Query().Get("code")
	if code == "" {
		writeError(w, http.StatusBadRequest, "code query parameter is required")
		return
	}
	resp, err := h.Clients.Ticket.GetTicketByCode(r.Context(), &ticketv1.GetTicketByCodeRequest{TicketCode: code})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	t := resp.GetTicket()
	if claims.Role != "admin" && t.GetUserId() != claims.UserID {
		writeError(w, http.StatusForbidden, "cannot verify another user's ticket")
		return
	}
	writeJSON(w, http.StatusOK, mapProtoTicket(t))
}

// CheckInTicket godoc
// @Summary      Check in attendee
// @Description  Marks ticket as checked in at venue entrance (admin)
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  CheckInRequest  true  "Ticket code"
// @Success      200   {object}  Ticket
// @Router       /api/v1/tickets/check-in [post]
func (h *Handler) CheckInTicket(w http.ResponseWriter, r *http.Request) {
	var req CheckInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.TicketCode == "" {
		writeError(w, http.StatusBadRequest, "ticketCode is required")
		return
	}
	resp, err := h.Clients.Ticket.CheckInTicket(r.Context(), &ticketv1.CheckInTicketRequest{TicketCode: req.TicketCode})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, mapProtoTicket(resp.GetTicket()))
}

// JoinWaitlist godoc
// @Summary      Join event waitlist
// @Description  Join waitlist when event is sold out
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  JoinWaitlistRequest  true  "Event to waitlist"
// @Success      201   {object}  WaitlistEntry
// @Router       /api/v1/waitlist [post]
func (h *Handler) JoinWaitlist(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	var req JoinWaitlistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.EventID == "" {
		writeError(w, http.StatusBadRequest, "eventId is required")
		return
	}
	resp, err := h.Clients.Ticket.JoinWaitlist(r.Context(), &ticketv1.JoinWaitlistRequest{
		UserId: claims.UserID, EventId: req.EventID,
	})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	e := resp.GetEntry()
	writeJSON(w, http.StatusCreated, WaitlistEntry{
		ID: e.GetId(), UserID: e.GetUserId(), EventID: e.GetEventId(), CreatedAt: e.GetCreatedAt(),
	})
}

// GetTicketByCode godoc
// @Summary      Get ticket by code
// @Description  Public lookup path used with ticket code in URL
// @Tags         tickets
// @Produce      json
// @Security     BearerAuth
// @Param        code  path  string  true  "Ticket code"
// @Success      200   {object}  Ticket
// @Router       /api/v1/tickets/code/{code} [get]
func (h *Handler) GetTicketByCode(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	code := mux.Vars(r)["code"]
	resp, err := h.Clients.Ticket.GetTicketByCode(r.Context(), &ticketv1.GetTicketByCodeRequest{TicketCode: code})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	t := resp.GetTicket()
	if claims.Role != "admin" && t.GetUserId() != claims.UserID {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}
	writeJSON(w, http.StatusOK, mapProtoTicket(t))
}
