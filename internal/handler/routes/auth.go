package routes

import (
	"net/http"
	"os"

	"github.com/garrychrstn/top-go/internal/repository"
	"github.com/garrychrstn/top-go/internal/util"
	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	auth repository.AuthRepository
}

func AuthInitHandler(auth repository.AuthRepository) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) Register(r chi.Router) {
	r.Post("/login", h.Login)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := util.DecodeJSON(w, r, &req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.auth.GetUser(r.Context(), req.Username)
	if err != nil {
		util.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err := util.VerifyPassword(user.Password, req.Password); err != nil {
		util.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := util.JWTGenerate(user.ID, "", false)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	secure := os.Getenv("ENV") == "production"
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Secure:   secure,
		Path:     "/",
		MaxAge:   86400,
		SameSite: http.SameSiteStrictMode,
	})

	util.WriteJSON(w, http.StatusOK, map[string]string{"message": "login successful"})
}
