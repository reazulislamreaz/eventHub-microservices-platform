package rest

import (
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterRoutes mounts documented REST API routes on the router.
func RegisterRoutes(r *mux.Router, h *Handler) {
	// Documentation
	r.PathPrefix("/api/docs/").Handler(http.StripPrefix("/api/docs/", APIDocsPage()))
	r.HandleFunc("/api/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api/docs/", http.StatusMovedPermanently)
	}).Methods("GET")
	r.HandleFunc("/api/v1/docs", APIDocs).Methods("GET")
	r.HandleFunc("/api/v1/graphql/schema", GraphQLSchema).Methods("GET")
	r.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api/docs/", http.StatusFound)
	}).Methods("GET")

	// Auth (public)
	r.HandleFunc("/api/v1/auth/register", h.Register).Methods("POST")
	r.HandleFunc("/api/v1/auth/login", h.Login).Methods("POST")

	// Users (public read)
	r.HandleFunc("/api/v1/users", h.ListUsers).Methods("GET")
	r.HandleFunc("/api/v1/users/{id}", h.GetUser).Methods("GET")

	// Events
	r.HandleFunc("/api/v1/events", h.ListEvents).Methods("GET")
	r.HandleFunc("/api/v1/events", h.RequireAdmin(h.CreateEvent)).Methods("POST")

	// Tickets (protected)
	r.HandleFunc("/api/v1/tickets", h.RequireAuth(h.BookTicket)).Methods("POST")
	r.HandleFunc("/api/v1/users/{id}/tickets", h.RequireAuth(h.GetUserTickets)).Methods("GET")
}
