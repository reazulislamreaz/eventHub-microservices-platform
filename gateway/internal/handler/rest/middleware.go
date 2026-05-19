package rest

import (
	"net/http"
	"strings"

	"github.com/eventhub/gateway/pkg/auth"
)

// RequireAuth validates Bearer JWT and attaches claims to request context.
func (h *Handler) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := bearerToken(r)
		if token == "" {
			writeError(w, http.StatusUnauthorized, "missing or invalid Authorization header")
			return
		}
		claims, err := h.JWT.Validate(token)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}
		next(w, r.WithContext(auth.WithClaims(r.Context(), claims)))
	}
}

// RequireAdmin requires authenticated user with admin role.
func (h *Handler) RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return h.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := auth.ClaimsFromContext(r.Context())
		if !ok || claims.Role != "admin" {
			writeError(w, http.StatusForbidden, "admin role required")
			return
		}
		next(w, r)
	})
}

func bearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return parts[1]
}
