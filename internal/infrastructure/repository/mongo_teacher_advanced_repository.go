package repository

import (
	"context"
	"errors"

	"github.com/srgjo27/e-learning/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoAssignmentRepository struct {
	collection *mongo.Collection
}

type MongoAssessmentRepository struct {
	collection *mongo.Collection
}

type MongoMessageRepository struct {
	collection *mongo.Collection
}

func NewMongoAssignmentRepository(c *mongo.Collection) *MongoAssignmentRepository {
	return &MongoAssignmentRepository{collection: c}
}

func NewMongoAssessmentRepository(c *mongo.Collection) *MongoAssessmentRepository {
	return &MongoAssessmentRepository{collection: c}
}

func NewMongoMessageRepository(c *mongo.Collection) *MongoMessageRepository {
	return &MongoMessageRepository{collection: c}
}

// --- Assignment ---
func (r *MongoAssignmentRepository) CreateAssignment(ctx context.Context, a *entity.Assignment) error {
	_, err := r.collection.InsertOne(ctx, a)
	return err
}

func (r *MongoAssignmentRepository) GetAssignment(ctx context.Context, id string) (*entity.Assignment, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var assignment entity.Assignment
	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&assignment)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("assignment not found")
		}
		return nil, err
	}
	return &assignment, nil
}

func (r *MongoAssignmentRepository) UpdateAssignment(ctx context.Context, a *entity.Assignment) error {
	if a.ID.IsZero() {
		return errors.New("assignment ID required")
	}
	filter := bson.M{"_id": a.ID}
	update := bson.M{
		"$set": bson.M{
			"title":       a.Title,
			"description": a.Description,
			"course_id":   a.CourseID,
			"due_date":    a.DueDate,
			"updated_at":  a.UpdatedAt,
		},
	}
	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("assignment not found")
	}
	return nil
}

func (r *MongoAssignmentRepository) DeleteAssignment(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("assignment not found")
	}
	return nil
}

func (r *MongoAssignmentRepository) ListAssignmentsByCourse(ctx context.Context, courseID primitive.ObjectID) ([]*entity.Assignment, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"course_id": courseID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var assignments []*entity.Assignment
	for cursor.Next(ctx) {
		var a entity.Assignment
		if err := cursor.Decode(&a); err != nil {
			return nil, err
		}
		assignments = append(assignments, &a)
	}
	return assignments, cursor.Err()
}

// --- Assessment ---
func (r *MongoAssessmentRepository) CreateAssessment(ctx context.Context, a *entity.Assessment) error {
	_, err := r.collection.InsertOne(ctx, a)
	return err
}

func (r *MongoAssessmentRepository) GetAssessment(ctx context.Context, id string) (*entity.Assessment, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var assessment entity.Assessment
	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&assessment)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("assessment not found")
		}
		return nil, err
	}
	return &assessment, nil
}

func (r *MongoAssessmentRepository) UpdateAssessment(ctx context.Context, a *entity.Assessment) error {
	if a.ID.IsZero() {
		return errors.New("assessment ID required")
	}
	filter := bson.M{"_id": a.ID}
	update := bson.M{
		"$set": bson.M{
			"title":       a.Title,
			"description": a.Description,
			"course_id":   a.CourseID,
			"date":        a.Date,
			"updated_at":  a.UpdatedAt,
		},
	}
	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("assessment not found")
	}
	return nil
}

func (r *MongoAssessmentRepository) DeleteAssessment(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("assessment not found")
	}
	return nil
}

func (r *MongoAssessmentRepository) ListAssessmentsByCourse(ctx context.Context, courseID primitive.ObjectID) ([]*entity.Assessment, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"course_id": courseID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var assessments []*entity.Assessment
	for cursor.Next(ctx) {
		var a entity.Assessment
		if err := cursor.Decode(&a); err != nil {
			return nil, err
		}
		assessments = append(assessments, &a)
	}
	return assessments, cursor.Err()
}

// --- Message ---
func (r *MongoMessageRepository) CreateMessage(ctx context.Context, m *entity.Message) error {
	_, err := r.collection.InsertOne(ctx, m)
	return err
}

func (r *MongoMessageRepository) GetMessage(ctx context.Context, id string) (*entity.Message, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var m entity.Message
	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&m)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("message not found")
		}
		return nil, err
	}
	return &m, nil
}

func (r *MongoMessageRepository) UpdateMessage(ctx context.Context, m *entity.Message) error {
	if m.ID.IsZero() {
		return errors.New("message ID required")
	}
	filter := bson.M{"_id": m.ID}
	update := bson.M{
		"$set": bson.M{
			"content":      m.Content,
			"receiver_ids": m.ReceiverIDs,
		},
	}
	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("message not found")
	}
	return nil
}

func (r *MongoMessageRepository) DeleteMessage(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("message not found")
	}
	return nil
}

func (r *MongoMessageRepository) ListMessagesBySender(ctx context.Context, senderID primitive.ObjectID) ([]*entity.Message, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"sender_id": senderID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []*entity.Message
	for cursor.Next(ctx) {
		var m entity.Message
		if err := cursor.Decode(&m); err != nil {
			return nil, err
		}
		messages = append(messages, &m)
	}
	return messages, cursor.Err()
}

func (r *MongoMessageRepository) ListMessagesForReceiver(ctx context.Context, receiverID primitive.ObjectID) ([]*entity.Message, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"receiver_ids": receiverID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []*entity.Message

	for cursor.Next(ctx) {
		var m entity.Message
		
		if err := cursor.Decode(&m); err != nil {
			return nil, err
		}
		messages = append(messages, &m)
	}

	return messages, cursor.Err()
}