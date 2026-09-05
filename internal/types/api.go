// Package types holds shared DTOs and request/response models used across
// the handler and repository layers.
package types

// ErrorResponse is the standard error body written by util.WriteError.
type ErrorResponse struct {
	Error string `json:"error"`
}
