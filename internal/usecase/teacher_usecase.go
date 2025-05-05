package usecase

import (
	"context"

	"github.com/srgjo27/e-learning/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TeacherUseCase struct {
	courseRepo CourseRepository
	classRepo  ClassRepository
	userRepo   UserRepository
}

func NewTeacherUseCase(courseRepo CourseRepository, classRepo ClassRepository, userRepo UserRepository) *TeacherUseCase {
	return &TeacherUseCase{
		courseRepo: courseRepo,
		classRepo:  classRepo,
		userRepo:   userRepo,
	}
}

func (t *TeacherUseCase) GetAssignedCourses(ctx context.Context, teacherID string) ([]*entity.Course, error) {
	oid, err := primitive.ObjectIDFromHex(teacherID)
	if err != nil {
		return nil, err
	}

	return t.courseRepo.ListCoursesByTeacher(ctx, oid)
}

func (t *TeacherUseCase) GetAssignedClasses(ctx context.Context, teacherID string) ([]*entity.Class, error) {
	oid, err := primitive.ObjectIDFromHex(teacherID)
	if err != nil {
		return nil, err
	}

	return t.classRepo.ListClassesByTeacher(ctx, oid)
}

func (t *TeacherUseCase) GetStudentsInClass(ctx context.Context, classID string) ([]*entity.User, error) {
	class, err := t.classRepo.GetClass(ctx, classID)
	if err != nil {
		return nil, err
	}
	if len(class.StudentIDs) == 0 {
		return []*entity.User{}, nil
	}

	return t.userRepo.FindUsersByIDs(ctx, class.StudentIDs)
}