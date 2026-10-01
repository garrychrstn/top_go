package routes

import (
	"net/http"
	"os"

	mw "github.com/garrychrstn/top-go/internal/handler/middleware"
	"github.com/garrychrstn/top-go/internal/repository"
	"github.com/garrychrstn/top-go/internal/types"
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
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.Login)
		r.Group(func(sub chi.Router) {
			sub.Use(mw.JWTAuth)
			sub.Post("/forgot-password", h.ForgotPassword)
		})
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req types.LoginRequest
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

	token, err := util.JWTGenerate(user.ID, user.Username)
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

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	claims := mw.GetJWTClaim(r.Context())
	if claims == nil {
		util.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req types.ForgotPassRequest
	if err := util.DecodeJSON(w, r, &req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Password == "" || req.OldPassword == "" {
		util.WriteError(w, http.StatusBadRequest, "both password and old password are required")
		return
	}

	user, err := h.auth.GetUser(r.Context(), claims.Username)
	if err != nil {
		util.WriteError(w, http.StatusUnauthorized, "user not found")
		return
	}

	if err := util.VerifyPassword(user.Password, req.OldPassword); err != nil {
		util.WriteError(w, http.StatusUnauthorized, "invalid old password")
		return
	}

	if err := h.auth.UpdatePassword(r.Context(), claims.UserID, req.Password); err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to update password")
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]string{"message": "password updated successfully"})
}
