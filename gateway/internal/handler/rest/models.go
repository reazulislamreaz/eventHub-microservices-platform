package rest

// ErrorResponse is returned for client and server errors.
// @Description API error payload
type ErrorResponse struct {
	Error   string `json:"error" example:"invalid request"`
	Code    int    `json:"code" example:"400"`
	Message string `json:"message,omitempty" example:"email is required"`
}

// User represents a platform user.
// @Description Registered user
type User struct {
	ID        string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email     string `json:"email" example:"alice@example.com"`
	Name      string `json:"name" example:"Alice"`
	Role      string `json:"role" example:"user" enums:"user,admin"`
	CreatedAt string `json:"createdAt" example:"2026-05-19T12:00:00Z"`
}

// Event represents a scheduled event.
// @Description Event with seat inventory
type Event struct {
	ID             string `json:"id" example:"550e8400-e29b-41d4-a716-446655440001"`
	Title          string `json:"title" example:"Go Conference 2026"`
	Description    string `json:"description" example:"Annual Go meetup"`
	Location       string `json:"location" example:"Dhaka, Bangladesh"`
	StartTime      string `json:"startTime" example:"2026-06-15T09:00:00Z"`
	EndTime        string `json:"endTime" example:"2026-06-15T18:00:00Z"`
	Capacity       int    `json:"capacity" example:"100"`
	AvailableSeats int    `json:"availableSeats" example:"95"`
	CreatedBy      string `json:"createdBy" example:"550e8400-e29b-41d4-a716-446655440000"`
	CreatedAt      string `json:"createdAt" example:"2026-05-19T12:00:00Z"`
}

// Ticket represents a booking for an event.
// @Description Confirmed event ticket
type Ticket struct {
	ID         string `json:"id" example:"550e8400-e29b-41d4-a716-446655440002"`
	UserID     string `json:"userId" example:"550e8400-e29b-41d4-a716-446655440000"`
	EventID    string `json:"eventId" example:"550e8400-e29b-41d4-a716-446655440001"`
	Status     string `json:"status" example:"confirmed"`
	TicketCode string `json:"ticketCode" example:"EH-a1b2c3d4e5f67890"`
	CreatedAt  string `json:"createdAt" example:"2026-05-19T12:30:00Z"`
}

// RegisterRequest registers a new user.
// @Description User registration body
type RegisterRequest struct {
	Email    string `json:"email" example:"alice@example.com" binding:"required"`
	Name     string `json:"name" example:"Alice" binding:"required"`
	Password string `json:"password" example:"SecurePass123" binding:"required" minLength:"8"`
}

// LoginRequest authenticates a user.
// @Description Login credentials
type LoginRequest struct {
	Email    string `json:"email" example:"admin@eventhub.io" binding:"required"`
	Password string `json:"password" example:"AdminPass123!" binding:"required"`
}

// AuthResponse contains JWT and user profile.
// @Description Authentication success response
type AuthResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIs..."`
	User  User   `json:"user"`
}

// CreateEventRequest creates a new event (admin).
// @Description Event creation body
type CreateEventRequest struct {
	Title       string `json:"title" example:"Go Conference 2026" binding:"required"`
	Description string `json:"description" example:"Annual Go community meetup"`
	Location    string `json:"location" example:"Dhaka, Bangladesh" binding:"required"`
	StartTime   string `json:"startTime" example:"2026-06-15T09:00:00Z" binding:"required"`
	EndTime     string `json:"endTime" example:"2026-06-15T18:00:00Z" binding:"required"`
	Capacity    int32  `json:"capacity" example:"100" binding:"required"`
}

// BookTicketRequest books a ticket for an event.
// @Description Ticket booking body
type BookTicketRequest struct {
	EventID string `json:"eventId" example:"550e8400-e29b-41d4-a716-446655440001" binding:"required"`
}

// APIDocsResponse describes available API surfaces.
// @Description API documentation index
type APIDocsResponse struct {
	Name        string            `json:"name" example:"EventHub Gateway API"`
	Version     string            `json:"version" example:"1.0"`
	Description string            `json:"description"`
	Endpoints   map[string]string `json:"endpoints"`
	GraphQL     GraphQLDocs       `json:"graphql"`
}

// GraphQLDocs points to GraphQL resources.
// @Description GraphQL API references
type GraphQLDocs struct {
	Playground string `json:"playground" example:"/"`
	Endpoint   string `json:"endpoint" example:"/query"`
	Schema     string `json:"schema" example:"/api/v1/graphql/schema"`
}
