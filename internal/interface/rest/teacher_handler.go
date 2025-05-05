package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/srgjo27/e-learning/internal/usecase"
)

type TeacherHandler struct {
	teacherUseCase *usecase.TeacherUseCase
}

func NewTeacherHandler(tu *usecase.TeacherUseCase) *TeacherHandler {
	return &TeacherHandler{
		teacherUseCase: tu,
	}
}

func (h * TeacherHandler) ListCourses(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	courses, err := h.teacherUseCase.GetAssignedCourses(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to get courses", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(courses)
}

func (h *TeacherHandler) ListClasses(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	classes, err := h.teacherUseCase.GetAssignedClasses(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to get classes", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(classes)
}

func (h *TeacherHandler) ListStudents(w http.ResponseWriter, r *http.Request) {
	classID := mux.Vars(r)["id"]
	if classID == "" {
		http.Error(w, "Class ID is required", http.StatusBadRequest)
		return
	}

	students, err := h.teacherUseCase.GetStudentsInClass(r.Context(), classID)
	if err != nil {
		http.Error(w, "Failed to get students", http.StatusInternalServerError)
		return
	}

	for _, s := range students {
		s.Password = ""
	}

	json.NewEncoder(w).Encode(students)
}
