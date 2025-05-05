package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/srgjo27/e-learning/internal/entity"
	"github.com/srgjo27/e-learning/internal/usecase"
)

type AdminHandler struct {
	authUseCase *usecase.AuthUseCase
}

func NewAdminHandler(u *usecase.AuthUseCase) *AdminHandler {
	return &AdminHandler{
		authUseCase: u,
	}
}

type updateRoleRequest struct {
	Role entity.Role `json:"role"`
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.authUseCase.ListAllUsers(r.Context())
	if err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	for _, u := range users {
		u.Password = ""
	}

	json.NewEncoder(w).Encode(users)
}

func (h *AdminHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/admin/users/")
	id = strings.TrimSuffix(id, "/role")
	if id == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	var req updateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Role != entity.RoleAdmin && req.Role != entity.RoleTeacher && req.Role != entity.RoleStudent {
		http.Error(w, "Invalid role", http.StatusBadRequest)
		return
	}

	err := h.authUseCase.UpdateUserRole(r.Context(), id, req.Role)
	if err != nil {
		http.Error(w, "Failed to update role", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "User role updated"})
}

func (h *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/admin/users/")
	if id == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	err := h.authUseCase.DeleteUser(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}