package rest

import (
	"encoding/json"
	"net/http"

	"github.com/srgjo27/e-learning/internal/usecase"
)

type TeacherAdvanceHandler struct {
	usecase *usecase.TeacherAdvancedUseCase
}

func NewTeacherAdvancedHandler(u *usecase.TeacherAdvancedUseCase) *TeacherAdvanceHandler {
	return &TeacherAdvanceHandler{
		usecase: u,
	}
}

func (h *TeacherAdvanceHandler) ListAssignments(w http.ResponseWriter, r *http.Request) {
	courseID := r.URL.Query().Get("course_id")
	if courseID == "" {
		http.Error(w, "course_id query param required", http.StatusBadRequest)
		return
	}

	assignments, err := h.usecase.ListAssessmentsByCourse(r.Context(), courseID)
	if err != nil {
		http.Error(w, "failed to list assignments", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(assignments)
}