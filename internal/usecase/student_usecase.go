package usecase

import (
	"context"

	"github.com/srgjo27/e-learning/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubmissionRepository interface {
	CreateSubmission(ctx context.Context, s *entity.Submission) error
	UpdateSubmission(ctx context.Context, s *entity.Submission) error
	GetSubmission(ctx context.Context, id string) (*entity.Submission, error)
	ListSubmissionsByStudent(ctx context.Context, studentID primitive.ObjectID) ([]*entity.Submission, error)
	ListSubmissionsByAssignmentAndStudent(ctx context.Context, assignmentID, studentID primitive.ObjectID) ([]*entity.Submission, error)
}

type StudentUseCase struct {
	courseRepo 		CourseRepository
	classRepo 		ClassRepository
	assignmentRepo 	AssignmentRepository
	assessmentRepo 	AssessmentRepository
	messageRepo 	MessageRepository
	submitRepo 		SubmissionRepository
	userRepo 		UserRepository
}

func NewStudentUseCase(
	courseRepo CourseRepository,
	classRepo ClassRepository,
	assignmentRepo AssignmentRepository,
	assessmentRepo AssessmentRepository,
	messageRepo MessageRepository,
	submitRepo SubmissionRepository,
	userRepo UserRepository,
) *StudentUseCase {
	return &StudentUseCase{
		courseRepo: courseRepo,
		classRepo: classRepo,
		assignmentRepo: assignmentRepo,
		assessmentRepo: assessmentRepo,
		messageRepo: messageRepo,
		submitRepo: submitRepo,
		userRepo: userRepo,
	}
}

func (s *StudentUseCase) GetEnrolledCourses(ctx context.Context, studentID string) ([]*entity.Course, error) {
	oid, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		return nil, err
	}

	classes, err := s.classRepo.ListClassesByStudent(ctx, oid)
	if err != nil {
		return nil, err
	}

	courseIDs := make(map[primitive.ObjectID]struct{})
	for _, cl := range classes {
		if !cl.CourseID.IsZero() {
			courseIDs[cl.CourseID] = struct{}{}
		}
	}

	var courses []*entity.Course

	for courseID := range courseIDs {
		course, err := s.courseRepo.GetCourse(ctx, courseID.Hex())
		if err == nil {
			courses = append(courses, course)
		}
	}

	return courses, nil
}

func (s *StudentUseCase) GetEnrolledClasses(ctx context.Context, studentID string) ([]*entity.Class, error) {
	oid, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		return nil, err
	}

	return s.classRepo.ListClassesByStudent(ctx, oid)
}

func (s *StudentUseCase) ListAssignmentsForStudent(ctx context.Context, studentID string) ([]*entity.Assignment, error) {
	// TODO:
	// For simplicity, we will return all assignments for now
	// In a real implementation, we would filter by the student's enrolled courses
	return s.assignmentRepo.ListAssignmentsByCourse(ctx, primitive.NewObjectID())
}

func (s *StudentUseCase) ListAssessmentsForStudent(ctx context.Context, studentID string) ([]*entity.Assessment, error) {
	return s.assessmentRepo.ListAssessments(ctx)
}

func (s *StudentUseCase) ListMessagesForStudent(ctx context.Context, studentID string) ([]*entity.Message, error) {
	oid, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		return nil, err
	}

	return s.messageRepo.ListMessagesForReceiver(ctx, oid)
}

func (s *StudentUseCase) SubmitAssignment(ctx context.Context, submission *entity.Submission) error {
	submission.SubmittedAt = submission.SubmittedAt.UTC()
	return s.submitRepo.CreateSubmission(ctx, submission)
}

func (s *StudentUseCase) UpdateSubmission(ctx context.Context, submission *entity.Submission) error {
	return s.submitRepo.UpdateSubmission(ctx, submission)
}

func (s *StudentUseCase) GetSubmission(ctx context.Context, submissionID string) (*entity.Submission, error) {
	return s.submitRepo.GetSubmission(ctx, submissionID)
}