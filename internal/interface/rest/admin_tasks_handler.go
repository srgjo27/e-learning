package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/srgjo27/e-learning/internal/entity"
	"github.com/srgjo27/e-learning/internal/usecase"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminTasksHandler struct {
	adminUseCase *usecase.AdminUseCase
}

func NewAdminTasksHandler(adminUseCase *usecase.AdminUseCase) *AdminTasksHandler {
	return &AdminTasksHandler{
		adminUseCase: adminUseCase,
	}
}

func(h *AdminTasksHandler) ListCourses(w http.ResponseWriter, r *http.Request) {
	courses, err := h.adminUseCase.ListCourses(r.Context())
	if err != nil {
		http.Error(w, "Failed to list courses", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(courses)
}

func(h *AdminTasksHandler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	var course entity.Course

	if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	course.CreatedAt = course.CreatedAt.UTC()
	if err := h.adminUseCase.CreateCourse(r.Context(), &course); err != nil {
		http.Error(w, "Failed to create course", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Course created successfully"})
}

func (h *AdminTasksHandler) GetCourse(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	course, err := h.adminUseCase.GetCourse(r.Context(), id)
	if err != nil {
		http.Error(w, "Course not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(course)
}

func (h *AdminTasksHandler) UpdateCourse(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var course entity.Course
	if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	course.ID = oid
	if err := h.adminUseCase.UpdateCourse(r.Context(), &course); err != nil {
		http.Error(w, "Failed to update course", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Course updated successfully"})
}

func (h *AdminTasksHandler) DeleteCourse(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.adminUseCase.DeleteCourse(r.Context(), id); err != nil {
		http.Error(w, "Failed to delete course", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AdminTasksHandler) ListClasses(w http.ResponseWriter, r *http.Request) {
	classes, err := h.adminUseCase.ListClasses(r.Context())
	if err != nil {
		http.Error(w, "Failed to list classes", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(classes)
}

func (h *AdminTasksHandler) CreateClass(w http.ResponseWriter, r *http.Request) {
	var class entity.Class
	if err := json.NewDecoder(r.Body).Decode(&class); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	class.CreatedAt = class.CreatedAt.UTC()
	if err := h.adminUseCase.CreateClass(r.Context(), &class); err != nil {
		http.Error(w, "Failed to create class", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Class created successfully"})
}

func (h *AdminTasksHandler) GetClass(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	class, err := h.adminUseCase.GetClass(r.Context(), id)
	if err != nil {
		http.Error(w, "Class not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(class)
}

func (h *AdminTasksHandler) UpdateClass(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var class entity.Class
	if err := json.NewDecoder(r.Body).Decode(&class); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	class.ID = oid
	if err := h.adminUseCase.UpdateClass(r.Context(), &class); err != nil {
		http.Error(w, "Failed to update class", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Class updated successfully"})
}

func (h *AdminTasksHandler) DeleteClass(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.adminUseCase.DeleteClass(r.Context(), id); err != nil {
		http.Error(w, "Failed to delete class", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AdminTasksHandler) ListAnnouncements(w http.ResponseWriter, r *http.Request) {
	anns, err := h.adminUseCase.ListAnnouncements(r.Context())
	if err != nil {
		http.Error(w, "Failed to list announcements", http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(anns)
}

func (h *AdminTasksHandler) CreateAnnouncement(w http.ResponseWriter, r *http.Request) {
	var ann entity.Announcement
	if err := json.NewDecoder(r.Body).Decode(&ann); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	ann.CreatedAt = ann.CreatedAt.UTC()
	if err := h.adminUseCase.CreateAnnouncement(r.Context(), &ann); err != nil {
		http.Error(w, "Failed to create announcement", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Announcement created successfully"})
}

func (h *AdminTasksHandler) GetAnnouncement(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	ann, err := h.adminUseCase.GetAnnouncement(r.Context(), id)
	if err != nil {
		http.Error(w, "Announcement not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(ann)
}

func (h *AdminTasksHandler) UpdateAnnouncement(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var ann entity.Announcement
	if err := json.NewDecoder(r.Body).Decode(&ann); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	ann.ID = oid
	if err := h.adminUseCase.UpdateAnnouncement(r.Context(), &ann); err != nil {
		http.Error(w, "Failed to update announcement", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Announcement updated successfully"})
}

func (h *AdminTasksHandler) DeleteAnnouncement(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.adminUseCase.DeleteAnnouncement(r.Context(), id); err != nil {
		http.Error(w, "Failed to delete announcement", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}