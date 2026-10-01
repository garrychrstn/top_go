package util

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/garrychrstn/top-go/internal/types"
)

// WriteJSON writes v as a JSON response with the given status code.
// A nil v writes an empty JSON body (204-style) but still sets the header.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}

// WriteError writes a JSON error body in the standard ErrorResponse shape.
func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, types.ErrorResponse{Error: msg})
}

// DecodeJSON decodes a JSON request body into dst, rejecting bodies over 1 MiB.
func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("invalid JSON body: %w", err)
	}
	return nil
}
