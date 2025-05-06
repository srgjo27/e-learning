package usecase

import (
	"context"

	"github.com/srgjo27/e-learning/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AssignmentRepository interface {
	CreateAssignment(ctx context.Context, a *entity.Assignment) error
	GetAssignment(ctx context.Context, id string) (*entity.Assignment, error)
	UpdateAssignment(ctx context.Context, a *entity.Assignment) error
	DeleteAssignment(ctx context.Context, id string) error
	ListAssignmentsByCourse(ctx context.Context, courseID primitive.ObjectID) ([]*entity.Assignment, error)
}

type AssessmentRepository interface {
	CreateAssessment(ctx context.Context, a *entity.Assessment) error
	GetAssessment(ctx context.Context, id string) (*entity.Assessment, error)
	UpdateAssessment(ctx context.Context, a *entity.Assessment) error
	DeleteAssessment(ctx context.Context, id string) error
	ListAssessmentsByCourse(ctx context.Context, courseID primitive.ObjectID) ([]*entity.Assessment, error)
}

type MessageRepository interface {
	CreateMessage(ctx context.Context, m *entity.Message) error
	GetMessage(ctx context.Context, id string) (*entity.Message, error)
	UpdateMessage(ctx context.Context, m *entity.Message) error
	DeleteMessage(ctx context.Context, id string) error
	ListMessagesBySender(ctx context.Context, senderID primitive.ObjectID) ([]*entity.Message, error)
}

type TeacherAdvancedUseCase struct {
	assignmentRepo AssignmentRepository
	assessmentRepo AssessmentRepository
	messageRepo    MessageRepository
}

func NewTeacherAdvancedUseCase(
	ar AssignmentRepository,
	asr AssessmentRepository,
	mr MessageRepository) *TeacherAdvancedUseCase {

	return &TeacherAdvancedUseCase{
		assignmentRepo: ar,
		assessmentRepo: asr,
		messageRepo:    mr,
	}
}

func (t *TeacherAdvancedUseCase) CreateAssignment(ctx context.Context, a *entity.Assignment) error {
	a.CreatedAt = a.CreatedAt.UTC()
	return t.assignmentRepo.CreateAssignment(ctx, a)
}

func (t *TeacherAdvancedUseCase) GetAssignment(ctx context.Context, id string) (*entity.Assignment, error) {
	return t.assignmentRepo.GetAssignment(ctx, id)
}

func (t *TeacherAdvancedUseCase) UpdateAssignment(ctx context.Context, a *entity.Assignment) error {
	return t.assignmentRepo.UpdateAssignment(ctx, a)
}

func (t *TeacherAdvancedUseCase) DeleteAssignment(ctx context.Context, id string) error {
	return t.assignmentRepo.DeleteAssignment(ctx, id)
}

func (t *TeacherAdvancedUseCase) ListAssignmentsByCourse(ctx context.Context, courseID string) ([]*entity.Assignment, error) {
	oid, err := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		return nil, err
	}
	return t.assignmentRepo.ListAssignmentsByCourse(ctx, oid)
}

func (t *TeacherAdvancedUseCase) CreateAssessment(ctx context.Context, a *entity.Assessment) error {
	a.CreatedAt = a.CreatedAt.UTC()
	return t.assessmentRepo.CreateAssessment(ctx, a)
}

func (t *TeacherAdvancedUseCase) GetAssessment(ctx context.Context, id string) (*entity.Assessment, error) {
	return t.assessmentRepo.GetAssessment(ctx, id)
}

func (t *TeacherAdvancedUseCase) UpdateAssessment(ctx context.Context, a *entity.Assessment) error {
	return t.assessmentRepo.UpdateAssessment(ctx, a)
}

func (t *TeacherAdvancedUseCase) DeleteAssessment(ctx context.Context, id string) error {
	return t.assessmentRepo.DeleteAssessment(ctx, id)
}

func (t *TeacherAdvancedUseCase) ListAssessmentsByCourse(ctx context.Context, courseID string) ([]*entity.Assessment, error) {
	oid, err := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		return nil, err
	}
	return t.assessmentRepo.ListAssessmentsByCourse(ctx, oid)
}

func (t *TeacherAdvancedUseCase) CreateMessage(ctx context.Context, m *entity.Message) error {
	m.CreatedAt = m.CreatedAt.UTC()
	return t.messageRepo.CreateMessage(ctx, m)
}

func (t *TeacherAdvancedUseCase) GetMessage(ctx context.Context, id string) (*entity.Message, error) {
	return t.messageRepo.GetMessage(ctx, id)
}

func (t *TeacherAdvancedUseCase) UpdateMessage(ctx context.Context, m *entity.Message) error {
	return t.messageRepo.UpdateMessage(ctx, m)
}

func (t *TeacherAdvancedUseCase) DeleteMessage(ctx context.Context, id string) error {
	return t.messageRepo.DeleteMessage(ctx, id)
}

func (t *TeacherAdvancedUseCase) ListMessagesBySender(ctx context.Context, senderID string) ([]*entity.Message, error) {
	oid, err := primitive.ObjectIDFromHex(senderID)
	if err != nil {
		return nil, err
	}
	return t.messageRepo.ListMessagesBySender(ctx, oid)
}