// internal/interfaces/http/handlers/auth.go
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/invoice-app-be/internal/domain/user"
	"github.com/invoice-app-be/internal/infrastructure/auth"
)

type AuthHandler struct {
	userService *user.Service
	jwtManager  *auth.JWTManager
}

func NewAuthHandler(userService *user.Service, jwtManager *auth.JWTManager) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtManager:  jwtManager,
	}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string  `json:"token"`
	User  UserDTO `json:"user"`
}

type UserDTO struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	authenticate, err := h.userService.Register(r.Context(), req.Email, req.Password, req.FullName)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.jwtManager.Generate(authenticate.ID, authenticate.Email)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	respondJSON(w, http.StatusCreated, AuthResponse{
		Token: token,
		User: UserDTO{
			ID:       authenticate.ID.String(),
			Email:    authenticate.Email,
			FullName: authenticate.FullName,
		},
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request")
		return
	}
	authenticate, err := h.userService.Authenticate(r.Context(), req.Email, req.Password)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Invalid cssssredentials")
		return
	}

	token, err := h.jwtManager.Generate(authenticate.ID, authenticate.Email)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	respondJSON(w, http.StatusOK, AuthResponse{
		Token: token,
		User: UserDTO{
			ID:       authenticate.ID.String(),
			Email:    authenticate.Email,
			FullName: authenticate.FullName,
		},
	})
}
