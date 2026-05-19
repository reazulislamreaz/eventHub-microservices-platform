package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/eventhub/gateway/pkg/auth"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func AuthMiddleware(jwtManager *auth.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if token := extractToken(r); token != "" {
				if claims, err := jwtManager.Validate(token); err == nil {
					ctx = auth.WithClaims(ctx, claims)
				}
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractToken(r *http.Request) string {
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

// DirectiveAuth enforces @auth directive on GraphQL fields.
func DirectiveAuth(ctx context.Context, _ interface{}, next graphql.Resolver) (interface{}, error) {
	if _, ok := auth.ClaimsFromContext(ctx); !ok {
		return nil, gqlerror.Errorf("authentication required")
	}
	return next(ctx)
}

// DirectiveHasRole enforces @hasRole directive.
func DirectiveHasRole(ctx context.Context, _ interface{}, next graphql.Resolver, role string) (interface{}, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, gqlerror.Errorf("authentication required")
	}
	if claims.Role != role {
		return nil, gqlerror.Errorf("insufficient permissions: requires role %s", role)
	}
	return next(ctx)
}
