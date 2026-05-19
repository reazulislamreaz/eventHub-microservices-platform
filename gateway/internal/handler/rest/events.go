package rest

import (
	"encoding/json"
	"net/http"

	eventv1 "github.com/eventhub/proto/gen/event/v1"
	"github.com/eventhub/gateway/pkg/auth"
)

// ListEvents godoc
// @Summary      List events
// @Description  Returns all published events with seat availability
// @Tags         events
// @Produce      json
// @Success      200  {array}   Event
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/events [get]
func (h *Handler) ListEvents(w http.ResponseWriter, r *http.Request) {
	resp, err := h.Clients.Event.ListEvents(r.Context(), &eventv1.ListEventsRequest{})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	events := make([]Event, 0, len(resp.GetEvents()))
	for _, e := range resp.GetEvents() {
		events = append(events, mapProtoEvent(e))
	}
	writeJSON(w, http.StatusOK, events)
}

// CreateEvent godoc
// @Summary      Create event
// @Description  Creates a new event (admin only). Requires Bearer JWT with admin role.
// @Tags         events
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      CreateEventRequest  true  "Event details"
// @Success      201   {object}  Event
// @Failure      400   {object}  ErrorResponse
// @Failure      401   {object}  ErrorResponse
// @Failure      403   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /api/v1/events [post]
func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	resp, err := h.Clients.Event.CreateEvent(r.Context(), &eventv1.CreateEventRequest{
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Capacity:    req.Capacity,
		CreatedBy:   claims.UserID,
	})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, mapProtoEvent(resp.GetEvent()))
}

func mapProtoEvent(e *eventv1.Event) Event {
	return Event{
		ID:             e.GetId(),
		Title:          e.GetTitle(),
		Description:    e.GetDescription(),
		Location:       e.GetLocation(),
		StartTime:      e.GetStartTime(),
		EndTime:        e.GetEndTime(),
		Capacity:       int(e.GetCapacity()),
		AvailableSeats: int(e.GetAvailableSeats()),
		CreatedBy:      e.GetCreatedBy(),
		CreatedAt:      e.GetCreatedAt(),
	}
}
