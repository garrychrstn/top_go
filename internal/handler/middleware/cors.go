package middleware

import (
	"net/http"
	"strconv"
	"strings"
)

// CORSConfig controls the CORS headers emitted for cross-origin requests.
type CORSConfig struct {
	AllowedOrigins   []string // "*" or explicit origins, e.g. ["http://localhost:5173"]
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int // seconds
}

// CORS adds CORS headers and short-circuits preflight (OPTIONS) requests.
func CORS(cfg CORSConfig) func(http.Handler) http.Handler {
	if len(cfg.AllowedMethods) == 0 {
		cfg.AllowedMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	}
	if len(cfg.AllowedHeaders) == 0 {
		cfg.AllowedHeaders = []string{"Content-Type", "Authorization", "X-Request-ID"}
	}
	if cfg.MaxAge == 0 {
		cfg.MaxAge = 3600
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" && originAllowed(cfg.AllowedOrigins, origin) {
				w.Header().Add("Vary", "Origin")

				// Credentials and "*" can't combine, so echo the specific origin.
				allowOrigin := origin
				if contains(cfg.AllowedOrigins, "*") {
					allowOrigin = "*"
				} else if cfg.AllowCredentials {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}
				w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.AllowedMethods, ", "))
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.AllowedHeaders, ", "))
				if len(cfg.ExposedHeaders) > 0 {
					w.Header().Set("Access-Control-Expose-Headers", strings.Join(cfg.ExposedHeaders, ", "))
				}
				w.Header().Set("Access-Control-Max-Age", strconv.Itoa(cfg.MaxAge))

				if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func originAllowed(allowed []string, origin string) bool {
	for _, a := range allowed {
		if a == "*" || strings.EqualFold(strings.TrimSuffix(a, "/"), strings.TrimSuffix(origin, "/")) {
			return true
		}
	}
	return false
}

func contains(items []string, target string) bool {
	for _, it := range items {
		if it == target {
			return true
		}
	}
	return false
}
