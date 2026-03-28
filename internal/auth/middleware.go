package auth

import (
	"net/http"

	"portlyn/internal/domain"
)

func (s *Service) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, groupIDs, err := s.AuthenticateRequest(r.Context(), r)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":{"code":"unauthorized","message":"missing or invalid session"}}`))
			return
		}

		ctx := ContextWithUser(r.Context(), user)
		ctx = ContextWithGroupIDs(ctx, groupIDs)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireRole(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := UserFromContext(r.Context())
			if !ok {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"error":{"code":"unauthorized","message":"missing auth context"}}`))
				return
			}
			if _, ok := allowed[user.Role]; !ok {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				_, _ = w.Write([]byte(`{"error":{"code":"forbidden","message":"insufficient permissions"}}`))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func CanRead(role string) bool {
	return role == domain.RoleAdmin || role == domain.RoleViewer
}
