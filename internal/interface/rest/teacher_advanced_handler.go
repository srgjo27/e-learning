package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/srgjo27/e-learning/internal/entity"
	"github.com/srgjo27/e-learning/internal/usecase"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TeacherAdvancedHandler struct {
	usecase *usecase.TeacherAdvancedUseCase
}

func NewTeacherAdvancedHandler(u *usecase.TeacherAdvancedUseCase) *TeacherAdvancedHandler {
	return &TeacherAdvancedHandler{
		usecase: u,
	}
}

func parseTimeISO8601(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

func (h *TeacherAdvancedHandler) ListAssignments(w http.ResponseWriter, r *http.Request) {
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

func (h *TeacherAdvancedHandler) CreateAssignment(w http.ResponseWriter, r *http.Request) {
	var a struct {
		Title	   	string `json:"title"`
		Description string `json:"description"`
		CourseID    string `json:"course_id"`
		DueDate		string `json:"due_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	dueDate, err := parseTimeISO8601(a.DueDate)
	if err != nil {
		http.Error(w, "invalid due_date", http.StatusBadRequest)
		return
	}

	newAssignment := &entity.Assignment{
		Title:       a.Title,
		Description: a.Description,
		CreatedAt:   time.Now(),
		DueDate:     dueDate,
	}

	id, err := primitive.ObjectIDFromHex(a.CourseID)
	if err != nil {
		http.Error(w, "invalid course_id", http.StatusBadRequest)
		return
	}

	newAssignment.CourseID = id

	if err := h.usecase.CreateAssignment(r.Context(), newAssignment); err != nil {
		http.Error(w, "failed to create assignment", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "assignment created successfully"})
}

func (h *TeacherAdvancedHandler) GetAssignment(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	assignment, err := h.usecase.GetAssignment(r.Context(), id)
	if err != nil {
		http.Error(w, "assignment not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(assignment)
}

func (h *TeacherAdvancedHandler) UpdateAssignment(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var a entity.Assessment
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	a.ID = oid
	if err := h.usecase.UpdateAssessment(r.Context(), &a); err != nil {
		http.Error(w, "Failed to update assessment", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "assignment updated successfully"})
}

func (h *TeacherAdvancedHandler) DeleteAssignment(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.usecase.DeleteAssignment(r.Context(), id); err != nil {
		http.Error(w, "Failed to delete assignment", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "assignment deleted successfully"})
}