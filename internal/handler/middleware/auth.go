package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/garrychrstn/top-go/internal/util"
)

type contextKey int

const jwtClaimKey contextKey = iota

// JWTAuth reads the token from the "token" cookie or Authorization header,
// parses it, and attaches the JWTClaim to the request context.
func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenStr string

		if cookie, err := r.Cookie("token"); err == nil {
			tokenStr = cookie.Value
		}

		if tokenStr == "" {
			authHeader := r.Header.Get("Authorization")
			if t, found := strings.CutPrefix(authHeader, "Bearer "); found {
				tokenStr = strings.TrimSpace(t)
			}
		}

		if tokenStr == "" {
			util.WriteError(w, http.StatusUnauthorized, "missing token")
			return
		}

		claims, err := util.JWTParse(tokenStr)
		if err != nil {
			util.WriteError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), jwtClaimKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetJWTClaim retrieves the JWTClaim from the context.
func GetJWTClaim(ctx context.Context) *util.JWTClaim {
	claims, _ := ctx.Value(jwtClaimKey).(*util.JWTClaim)
	return claims
}
