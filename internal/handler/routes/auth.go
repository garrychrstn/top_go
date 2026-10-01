package routes

import "github.com/garrychrstn/top-go/internal/repository"

type AuthHandler struct {
	auth repository.AuthRepository
}

func AuthInitHandler(auth repository.AuthRepository) *AuthHandler {
	return &AuthHandler{auth: auth}
}
