package middleware

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/garrychrstn/top-go/internal/util"
)

// RequireBearerToken guards routes with a static bearer token supplied via
// config (e.g. AUTH_BEARER_TOKEN). An empty expected value disables the check
// so unconfigured environments still work in local dev.
func RequireBearerToken(expected string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if expected == "" {
				next.ServeHTTP(w, r)
				return
			}

			auth := r.Header.Get("Authorization")
			token, ok := strings.CutPrefix(auth, "Bearer ")
			if !ok || subtle.ConstantTimeCompare([]byte(token), []byte(expected)) != 1 {
				util.WriteError(w, http.StatusUnauthorized, "unauthorized")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
