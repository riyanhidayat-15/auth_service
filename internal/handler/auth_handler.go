package handler

import (
	"auth-service/internal/services"
	"encoding/json"
	"errors"
	"net/http"
)

type AuthHandler struct{
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{ authService: authService}
}

type credentialRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req credentialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "body tidak valid"})
		return
	}

	if req.Username == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "username atau password salah"})
		return
	}

	user, err := h.authService.Register(r.Context(), req.Username, req.Password)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"id": user.ID,
		"username":user.Username,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req credentialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "body tidak valid"})
		return
	}

	token, err := h.authService.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "terjadi kesalahan server"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}