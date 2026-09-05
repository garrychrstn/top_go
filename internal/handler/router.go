package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	mw "github.com/garrychrstn/go-v1/internal/handler/middleware"
)

// NewRouter assembles the full HTTP stack: middleware chain + routes.
func NewRouter(users *UserHandler, corsOrigins []string) http.Handler {
	r := chi.NewRouter()

	r.Use(mw.Recoverer)
	r.Use(mw.RequestID)
	r.Use(mw.Logger)
	r.Use(mw.CORS(mw.CORSConfig{AllowedOrigins: corsOrigins}))

	r.Get("/health", Health)

	users.Register(r)

	return r
}
