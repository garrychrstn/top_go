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

type ItemHandler struct {
	items repository.ItemRepository
}

func ItemInitHandler(items repository.ItemRepository) *ItemHandler {
	return &ItemHandler{items: items}
}

func (h *ItemHandler) Register(r chi.Router) {
	r.Group(func(sub chi.Router) {
		sub.Use(mw.JWTAuth)
		sub.Get("/items", h.List)
		sub.Post("/items", h.Create)
		sub.Put("/items/{id}", h.Update)
		sub.Delete("/items/{id}", h.Delete)
	})
}

type ItemRequest struct {
	Name  string `json:"name"`
	Price string `json:"price"`
}

func (h *ItemHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.items.List(r.Context())
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to list items")
		return
	}
	util.WriteJSON(w, http.StatusOK, items)
}

func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req ItemRequest
	if err := util.DecodeJSON(w, r, &req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var numericPrice pgtype.Numeric
	if err := numericPrice.Scan(req.Price); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid price format")
		return
	}

	item, err := h.items.Create(r.Context(), req.Name, numericPrice)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to create item")
		return
	}

	util.WriteJSON(w, http.StatusCreated, item)
}

func (h *ItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	var req ItemRequest
	if err := util.DecodeJSON(w, r, &req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var numericPrice pgtype.Numeric
	if err := numericPrice.Scan(req.Price); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid price format")
		return
	}

	item, err := h.items.Update(r.Context(), id, req.Name, numericPrice)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to update item")
		return
	}

	util.WriteJSON(w, http.StatusOK, item)
}

func (h *ItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	if err := h.items.Delete(r.Context(), id); err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to delete item")
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]string{"message": "item deleted successfully"})
}
