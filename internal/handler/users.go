package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/garrychrstn/top-go/internal/repository"
	"github.com/garrychrstn/top-go/internal/types"
	"github.com/garrychrstn/top-go/internal/util"
)

// UserHandler handles /users requests and calls the repository directly
// (no service layer in v1).
type UserHandler struct {
	users repository.UserRepository
}

// NewUserHandler injects the repository dependency.
func NewUserHandler(users repository.UserRepository) *UserHandler {
	return &UserHandler{users: users}
}

// Register mounts the /users routes on a router.
func (h *UserHandler) Register(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Get("/{userID}", h.GetByID)
	})
}

// List handles GET /users.
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.users.List(r.Context())
	if err != nil {
		slog.Error("list users", "error", err)
		util.WriteError(w, http.StatusInternalServerError, "failed to list users")
		return
	}
	util.WriteJSON(w, http.StatusOK, users)
}

// Create handles POST /users.
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req types.CreateUserRequest
	if err := util.DecodeJSON(w, r, &req); err != nil {
		util.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	if req.Name == "" {
		util.WriteError(w, http.StatusBadRequest, "name is required")
		return
	}
	addr, err := mail.ParseAddress(req.Email)
	if err != nil || addr.Address != req.Email {
		util.WriteError(w, http.StatusBadRequest, "email is invalid")
		return
	}

	user, err := h.users.Create(r.Context(), req)
	if err != nil {
		if errors.Is(err, repository.ErrEmailTaken) {
			util.WriteError(w, http.StatusConflict, "email is already registered")
			return
		}
		slog.Error("create user", "error", err)
		util.WriteError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	util.WriteJSON(w, http.StatusCreated, user)
}

// GetByID handles GET /users/{userID}.
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.users.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			util.WriteError(w, http.StatusNotFound, "user not found")
			return
		}
		slog.Error("get user by id", "error", err)
		util.WriteError(w, http.StatusInternalServerError, "failed to get user")
		return
	}

	util.WriteJSON(w, http.StatusOK, user)
}
