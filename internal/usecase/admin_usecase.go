package usecase

import (
	"context"

	"github.com/srgjo27/e-learning/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CourseRepository interface {
	CreateCourse(ctx context.Context, course *entity.Course) error
	GetCourse(ctx context.Context, id string) (*entity.Course, error)
	UpdateCourse(ctx context.Context, course *entity.Course) error 
	DeleteCourse(ctx context.Context, id string) error
	ListCourses(ctx context.Context) ([]*entity.Course, error)
	ListCoursesByTeacher(ctx context.Context, teacherID primitive.ObjectID) ([]*entity.Course, error)
}

type ClassRepository interface {
	CreateClass(ctx context.Context, class *entity.Class) error
	GetClass(ctx context.Context, id string) (*entity.Class, error)
	UpdateClass(ctx context.Context, class *entity.Class) error
	DeleteClass(ctx context.Context, id string) error
	ListClasses(ctx context.Context) ([]*entity.Class, error)
	ListClassesByTeacher(ctx context.Context, teacherID primitive.ObjectID) ([]*entity.Class, error)
	ListClassesByStudent(ctx context.Context, studentID primitive.ObjectID) ([]*entity.Class, error)
}

type AnnouncementRepository interface {
	CreateAnnouncement(ctx context.Context, ann *entity.Announcement) error
	GetAnnouncement(ctx context.Context, id string) (*entity.Announcement, error)
	UpdateAnnouncement(ctx context.Context, ann *entity.Announcement) error
	DeleteAnnouncement(ctx context.Context, id string) error
	ListAnnouncements(ctx context.Context) ([]*entity.Announcement, error)
}

type AdminUseCase struct {
	courseRepo 	 	 CourseRepository
	classRepo 	 	 ClassRepository
	announcementRepo AnnouncementRepository
}

func NewAdminUseCase(courseRepo CourseRepository, classRepo ClassRepository, announcementRepo AnnouncementRepository) *AdminUseCase {
	return &AdminUseCase{
		courseRepo: 	  courseRepo,
		classRepo: 		  classRepo,
		announcementRepo: announcementRepo,
	}
}

// --- Course ---
func (a *AdminUseCase) CreateCourse(ctx context.Context, course *entity.Course) error {
	return a.courseRepo.CreateCourse(ctx, course)
}

func (a *AdminUseCase) GetCourse(ctx context.Context, id string) (*entity.Course,error) {
	return a.courseRepo.GetCourse(ctx, id)
}

func (a *AdminUseCase) UpdateCourse(ctx context.Context, course *entity.Course) error {
	return a.courseRepo.UpdateCourse(ctx, course)
}

func (a *AdminUseCase) DeleteCourse(ctx context.Context, id string) error {
	return a.courseRepo.DeleteCourse(ctx, id)
}

func (a *AdminUseCase) ListCourses(ctx context.Context) ([]*entity.Course, error) {
	return a.courseRepo.ListCourses(ctx)
}

// --- Class ---
func (a *AdminUseCase) CreateClass(ctx context.Context, class *entity.Class) error {
	return a.classRepo.CreateClass(ctx, class)
}

func (a *AdminUseCase) GetClass(ctx context.Context, id string) (*entity.Class, error) {
	return a.classRepo.GetClass(ctx, id)
}

func (a *AdminUseCase) UpdateClass(ctx context.Context, class *entity.Class) error {
	return a.classRepo.UpdateClass(ctx, class)
}

func (a *AdminUseCase) DeleteClass(ctx context.Context, id string) error {
	return a.classRepo.DeleteClass(ctx, id)
}

func (a *AdminUseCase) ListClasses(ctx context.Context) ([]*entity.Class, error) {
	return a.classRepo.ListClasses(ctx)
}

// --- Announcement ---
func (a *AdminUseCase) CreateAnnouncement(ctx context.Context, ann *entity.Announcement) error {
	return a.announcementRepo.CreateAnnouncement(ctx, ann)
}

func (a *AdminUseCase) GetAnnouncement(ctx context.Context, id string) (*entity.Announcement, error) {
	return a.announcementRepo.GetAnnouncement(ctx, id)
}

func (a *AdminUseCase) UpdateAnnouncement(ctx context.Context, ann *entity.Announcement) error {
	return a.announcementRepo.UpdateAnnouncement(ctx, ann)
}

func (a *AdminUseCase) DeleteAnnouncement(ctx context.Context, id string) error {
	return a.announcementRepo.DeleteAnnouncement(ctx, id)
}

func (a *AdminUseCase) ListAnnouncements(ctx context.Context) ([]*entity.Announcement, error) {
	return a.announcementRepo.ListAnnouncements(ctx)
}