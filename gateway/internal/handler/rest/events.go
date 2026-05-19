package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	eventv1 "github.com/eventhub/proto/gen/event/v1"
	"github.com/eventhub/gateway/pkg/auth"
	"github.com/gorilla/mux"
)

// ListEvents godoc
// @Summary      List events (paginated)
// @Description  Search and filter published events with pagination
// @Tags         events
// @Produce      json
// @Param        page      query  int     false  "Page number"       default(1)
// @Param        pageSize  query  int     false  "Items per page"    default(20)
// @Param        search    query  string  false  "Search title/description"
// @Param        location  query  string  false  "Filter by location"
// @Param        status    query  string  false  "Event status (admin)"
// @Param        category  query  string  false  "Category filter (music, tech, sports, ...)"
// @Success      200       {object}  EventPage
// @Failure      500       {object}  ErrorResponse
// @Router       /api/v1/events [get]
func (h *Handler) ListEvents(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("pageSize"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	resp, err := h.Clients.Event.ListEvents(r.Context(), &eventv1.ListEventsRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   q.Get("search"),
		Location: q.Get("location"),
		Status:   q.Get("status"),
		Category: q.Get("category"),
	})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	events := make([]Event, 0, len(resp.GetEvents()))
	for _, e := range resp.GetEvents() {
		events = append(events, mapProtoEvent(e))
	}
	writeJSON(w, http.StatusOK, EventPage{
		Events:   events,
		Total:    int(resp.GetTotal()),
		Page:     int(resp.GetPage()),
		PageSize: int(resp.GetPageSize()),
	})
}

// GetEvent godoc
// @Summary      Get event by ID
// @Description  Returns a single event
// @Tags         events
// @Produce      json
// @Param        id   path      string  true  "Event ID"
// @Success      200  {object}  Event
// @Failure      404  {object}  ErrorResponse
// @Router       /api/v1/events/{id} [get]
func (h *Handler) GetEvent(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	resp, err := h.Clients.Event.GetEvent(r.Context(), &eventv1.GetEventRequest{Id: id})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, mapProtoEvent(resp.GetEvent()))
}

// CreateEvent godoc
// @Summary      Create event
// @Description  Creates a new event (admin only)
// @Tags         events
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      CreateEventRequest  true  "Event details"
// @Success      201   {object}  Event
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
		Title: req.Title, Description: req.Description, Location: req.Location,
		Category: req.Category, PriceCents: req.PriceCents,
		StartTime: req.StartTime, EndTime: req.EndTime, Capacity: req.Capacity,
		CreatedBy: claims.UserID,
	})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, mapProtoEvent(resp.GetEvent()))
}

// CancelEvent godoc
// @Summary      Cancel event
// @Description  Marks event as cancelled (admin only); blocks new bookings
// @Tags         events
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Event ID"
// @Success      200  {object}  Event
// @Router       /api/v1/events/{id}/cancel [post]
func (h *Handler) CancelEvent(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	resp, err := h.Clients.Event.CancelEvent(r.Context(), &eventv1.CancelEventRequest{Id: id})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, mapProtoEvent(resp.GetEvent()))
}

func mapProtoEvent(e *eventv1.Event) Event {
	return Event{
		ID: e.GetId(), Title: e.GetTitle(), Description: e.GetDescription(),
		Location: e.GetLocation(), Category: e.GetCategory(), PriceCents: e.GetPriceCents(),
		StartTime: e.GetStartTime(), EndTime: e.GetEndTime(),
		Capacity: int(e.GetCapacity()), AvailableSeats: int(e.GetAvailableSeats()),
		Status: e.GetStatus(), CreatedBy: e.GetCreatedBy(), CreatedAt: e.GetCreatedAt(),
	}
}
