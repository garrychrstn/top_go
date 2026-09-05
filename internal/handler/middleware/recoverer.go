// Package middleware provides the HTTP middleware stack for the router:
// logging, panic recovery, CORS, request IDs, and bearer-token auth.
package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/garrychrstn/go-v1/internal/util"
)

// Recoverer catches panics in downstream handlers, logs the stack, and
// writes a generic 500 instead of killing the server.
func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("panic recovered",
					"panic", rec,
					"request_id", GetRequestID(r.Context()),
					"stack", string(debug.Stack()),
				)
				util.WriteError(w, http.StatusInternalServerError, "internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
