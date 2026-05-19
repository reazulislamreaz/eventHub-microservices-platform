package main

import (
	"os"

	"github.com/eventhub/gateway/docs"
)

// configureSwaggerHost sets the OpenAPI host Swagger UI uses for "Try it out".
// Empty host = same origin as the browser (fixes CORS when Docker maps 8082→8080).
func configureSwaggerHost(swaggerHost string) {
	if swaggerHost != "" {
		docs.SwaggerInfo.Host = swaggerHost
		return
	}
	if h := os.Getenv("SWAGGER_HOST"); h != "" {
		docs.SwaggerInfo.Host = h
		return
	}
	// Relative requests: Swagger UI uses window.location.host (e.g. localhost:8082).
	docs.SwaggerInfo.Host = ""
}
