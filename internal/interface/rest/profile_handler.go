package rest

import (
	"encoding/json"
	"net/http"

	"github.com/srgjo27/e-learning/internal/entity"
	"github.com/srgjo27/e-learning/internal/usecase"
)

type ProfileHandler struct {
	authUseCase *usecase.AuthUseCase
}

func NewProfileHandler(u *usecase.AuthUseCase) *ProfileHandler {
	return &ProfileHandler{
		authUseCase : u,
	}
}

type profileUpdateRequest struct {
	Email string `json:"email,omitempty"`
	Password string `json:"new_password,omitempty"`
}

func (h *ProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID")
	if userID == nil {
		http.Error(w, "Unouthorized", http.StatusUnauthorized)
		return
	}

	uidStr, ok := userID.(string)
	if !ok {
		http.Error(w, "Unouthorized", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getProfile(w, r, uidStr)
	case http.MethodPut:
		h.updateProfile(w, r, uidStr)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ProfileHandler) getProfile(w http.ResponseWriter, r *http.Request, userID string) {
	user, err := h.authUseCase.GetProfile(r.Context(), userID)
	if err != nil {
		if err == entity.ErrUserNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get profile", http.StatusInternalServerError)
		return
	}

	user.Password = ""
	json.NewEncoder(w).Encode(user)
}

func (h *ProfileHandler) updateProfile(w http.ResponseWriter, r *http.Request, userID string) {
	var req profileUpdateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := h.authUseCase.UpdateProfile(r.Context(), userID, req.Email, req.Password)
	if err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated"})
}