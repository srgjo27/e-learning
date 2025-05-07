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

// --- Assignments ---

func (h *TeacherAdvancedHandler) ListAssignments(w http.ResponseWriter, r *http.Request) {
	courseID := r.URL.Query().Get("course_id")
	if courseID == "" {
		http.Error(w, "course_id query param required", http.StatusBadRequest)
		return
	}

	assignments, err := h.usecase.ListAssignmentsByCourse(r.Context(), courseID)
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
	var a entity.Assignment
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
	if err := h.usecase.UpdateAssignment(r.Context(), &a); err != nil {
		http.Error(w, "Failed to update assignment", http.StatusInternalServerError)
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

// --- Assessments ---

func (h *TeacherAdvancedHandler) ListAssessments(w http.ResponseWriter, r *http.Request) {
	courseID := r.URL.Query().Get("course_id")
	if courseID == "" {
		http.Error(w, "course_id query param required", http.StatusBadRequest)
		return
	}

	assessments, err := h.usecase.ListAssessmentsByCourse(r.Context(), courseID)
	if err != nil {
		http.Error(w, "failed to list assessments", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(assessments)
}

func (h *TeacherAdvancedHandler) CreateAssessment(w http.ResponseWriter, r *http.Request) {
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

	newAssessment := &entity.Assessment{
		Title:       a.Title,
		Description: a.Description,
		CreatedAt:   time.Now(),
		Date:     	 dueDate,
	}

	id, err := primitive.ObjectIDFromHex(a.CourseID)
	if err != nil {
		http.Error(w, "invalid course_id", http.StatusBadRequest)
		return
	}

	newAssessment.CourseID = id

	if err := h.usecase.CreateAssessment(r.Context(), newAssessment); err != nil {
		http.Error(w, "failed to create assessment", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "assessment created successfully"})
}

func (h *TeacherAdvancedHandler) GetAssessment(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	assignment, err := h.usecase.GetAssessment(r.Context(), id)
	if err != nil {
		http.Error(w, "assignment not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(assignment)
}

func (h *TeacherAdvancedHandler) UpdateAssessment(w http.ResponseWriter, r *http.Request) {
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

	json.NewEncoder(w).Encode(map[string]string{"message": "assessment updated successfully"})
}

func (h *TeacherAdvancedHandler) DeleteAssessment(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.usecase.DeleteAssessment(r.Context(), id); err != nil {
		http.Error(w, "Failed to delete assessment", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "assessment deleted successfully"})
}

// --- Messages ---

func (h *TeacherAdvancedHandler) ListMessages(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	messages, err := h.usecase.ListMessagesBySender(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to list messages", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(messages)
}

func (h *TeacherAdvancedHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var m struct {
		Content   	string 	 `json:"content"`
		ReceiverIDs []string `json:"receiver_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	if m.Content == "" || len(m.ReceiverIDs) == 0 {
		http.Error(w, "content and receiver_ids required", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	senderOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "invalid sender", http.StatusBadRequest)
		return
	}

	var receivers []primitive.ObjectID
	for _, rid := range m.ReceiverIDs {
		oid, err := primitive.ObjectIDFromHex(rid)
		if err != nil {
			http.Error(w, "invalid receiver ID", http.StatusBadRequest)
			return
		}
		receivers = append(receivers, oid)
	}

	message := &entity.Message{
		SenderID:    senderOID,
		ReceiverIDs: receivers,
		Content:     m.Content,
		CreatedAt:   time.Now(),
	}

	if err := h.usecase.CreateMessage(r.Context(), message); err != nil {
		http.Error(w, "failed to create message", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "message created successfully"})
}