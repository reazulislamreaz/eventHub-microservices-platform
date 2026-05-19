package rest

import (
	_ "embed"
	"net/http"
)

//go:embed schema.graphql
var graphQLSchema string

// APIDocs godoc
// @Summary      API documentation index
// @Description  Returns links to Swagger UI, GraphQL Playground, and schema endpoints
// @Tags         documentation
// @Produce      json
// @Success      200  {object}  APIDocsResponse
// @Router       /api/v1/docs [get]
func APIDocs(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, APIDocsResponse{
		Name:        "EventHub Gateway API",
		Version:     "1.0",
		Description: "Microservices gateway exposing REST (Swagger) and GraphQL APIs for event management, registration, and ticketing.",
		Endpoints: map[string]string{
			"documentation": "/api/docs/",
			"swagger":       "/swagger/index.html",
			"openapi":         "/swagger/doc.json",
			"health":          "/health",
			"readiness":       "/ready",
			"register":        "POST /api/v1/auth/register",
			"login":           "POST /api/v1/auth/login",
			"users":           "GET /api/v1/users",
			"events":          "GET /api/v1/events",
			"bookTicket":      "POST /api/v1/tickets",
		},
		GraphQL: GraphQLDocs{
			Playground: "/",
			Endpoint:   "/query",
			Schema:     "/api/v1/graphql/schema",
		},
	})
}

// GraphQLSchema godoc
// @Summary      GraphQL schema (SDL)
// @Description  Returns the GraphQL schema definition language file used by the gateway
// @Tags         documentation
// @Produce      plain
// @Success      200  {string}  string  "GraphQL SDL"
// @Router       /api/v1/graphql/schema [get]
func GraphQLSchema(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(graphQLSchema))
}

