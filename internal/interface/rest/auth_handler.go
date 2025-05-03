package rest

import (
	"encoding/json"
	"net/http"

	"github.com/srgjo27/e-learning/internal/entity"
	"github.com/srgjo27/e-learning/internal/usecase"
)

type AuthHandler struct {
	authUseCase *usecase.AuthUseCase
}

func NewAuthHandler(u *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: u,
	}
}

type registerRequest struct {
	Email	 string      `json:"email"`
	Password string      `json:"password"`
	Role 	 entity.Role `json:"role"`
}

type loginRequest struct {
	Email	 string `json:"email"`
	Password string `json:"password"`
}

type passwordResetRequest struct {
	Email string `json:"email"`
}

type passwordResetResetRequest struct {
	Token   	string `json:"token"`
	NewPassword string `json:"new_password"`
}

func (h *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req registerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password required", http.StatusBadRequest)
		return
	}

	err := h.authUseCase.Register(r.Context(), req.Email, req.Password, req.Role)
	if err != nil {
		if err == entity.ErrEmailExists {
			http.Error(w, "Email already registered", http.StatusConflict)
			return
		}
		http.Error(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req loginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password required", http.StatusBadRequest)
		return
	}

	token, err := h.authUseCase.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if err == entity.ErrUserNotFound || err == entity.ErrInvalidPassword{
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Login failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *AuthHandler) HandlePasswordResetRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req passwordResetRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	
	if req.Email == "" {
		http.Error(w, "Email required", http.StatusBadRequest)
		return
	}

	token, err := h.authUseCase.RequestPasswordReset(r.Context(), req.Email)
	if err != nil {
		if err == entity.ErrUserNotFound {
			http.Error(w, "Email not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Could not generate reset token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"reset_token": token})
}

func (h *AuthHandler) HandlePasswordReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req passwordResetResetRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Token == "" || req.NewPassword == "" {
		http.Error(w, "Token and new_password required", http.StatusBadRequest)
		return
	}

	err := h.authUseCase.ResetPassword(r.Context(), req.Token, req.NewPassword)
	if err != nil {
		if err == entity.ErrInvalidToken {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Password reset failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Password reset successful"})
}