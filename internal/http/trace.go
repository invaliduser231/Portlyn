package http

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) traceContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetReqID(r.Context())
		if requestID != "" {
			r.Header.Set("X-Request-Id", requestID)
			w.Header().Set("X-Request-Id", requestID)
			w.Header().Set("X-Trace-Id", requestID)
		}
		next.ServeHTTP(w, r)
	})
}
