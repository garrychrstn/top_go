package routes

import (
	"net/http"

	mw "github.com/garrychrstn/top-go/internal/handler/middleware"
	"github.com/garrychrstn/top-go/internal/repository"
	"github.com/garrychrstn/top-go/internal/util"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CustomerHandler struct {
	customers repository.CustomerRepository
}

func CustomerInitHandler(customers repository.CustomerRepository) *CustomerHandler {
	return &CustomerHandler{customers: customers}
}

func (h *CustomerHandler) Register(r chi.Router) {
	r.Group(func(sub chi.Router) {
		sub.Use(mw.JWTAuth)
		sub.Get("/customers", h.List)
		sub.Post("/customers", h.Create)
		sub.Get("/customers/{id}", h.GetByID)
		sub.Put("/customers/{id}", h.Update)
		sub.Delete("/customers/{id}", h.Delete)
	})
}

type CustomerRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
}

func (h *CustomerHandler) List(w http.ResponseWriter, r *http.Request) {
	customers, err := h.customers.List(r.Context())
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to list customers")
		return
	}
	util.WriteJSON(w, http.StatusOK, customers)
}

func (h *CustomerHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid customer id")
		return
	}

	c, err := h.customers.GetByID(r.Context(), id)
	if err != nil {
		util.WriteError(w, http.StatusNotFound, "customer not found")
		return
	}
	util.WriteJSON(w, http.StatusOK, c)
}

func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CustomerRequest
	if err := util.DecodeJSON(w, r, &req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var addr pgtype.Text
	if req.Address != "" {
		addr = pgtype.Text{String: req.Address, Valid: true}
	}

	c, err := h.customers.Create(r.Context(), req.Name, req.PhoneNumber, addr)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to create customer")
		return
	}

	util.WriteJSON(w, http.StatusCreated, c)
}

func (h *CustomerHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid customer id")
		return
	}

	var req CustomerRequest
	if err := util.DecodeJSON(w, r, &req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var addr pgtype.Text
	if req.Address != "" {
		addr = pgtype.Text{String: req.Address, Valid: true}
	}

	c, err := h.customers.Update(r.Context(), id, req.Name, req.PhoneNumber, addr)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to update customer")
		return
	}

	util.WriteJSON(w, http.StatusOK, c)
}

func (h *CustomerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid customer id")
		return
	}

	if err := h.customers.Delete(r.Context(), id); err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to delete customer")
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]string{"message": "customer deleted successfully"})
}
