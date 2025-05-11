package repository

import (
	"context"
	"errors"

	"github.com/srgjo27/e-learning/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoSubmissionRepository struct {
	collection *mongo.Collection
}

func NewMongoSubmissionRepository(c *mongo.Collection) *MongoSubmissionRepository {
	return &MongoSubmissionRepository{collection: c}
}

func (r *MongoSubmissionRepository) CreateSubmission(ctx context.Context, s *entity.Submission) error {
	_, err := r.collection.InsertOne(ctx, s)
	return err
}

func (r *MongoSubmissionRepository) UpdateSubmission(ctx context.Context, s *entity.Submission) error {
	if s.ID.IsZero() {
		return errors.New("Submission ID required")
	}

	filter := bson.M{"_id": s.ID}
	update := bson.M{
		"$set": bson.M{
			"content": 		s.Content,
			"grade":  		s.Grade,
			"feedback": 	s.Feedback,
			"submitted_at": s.SubmittedAt,
		},
	}
	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return entity.ErrSubmissionNotFound
	}

	return nil
}

func(r *MongoSubmissionRepository) GetSubmission(ctx context.Context, id string) (*entity.Submission, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var s entity.Submission
	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&s)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, entity.ErrSubmissionNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *MongoSubmissionRepository) ListSubmissionsByStudent(ctx context.Context, studentID primitive.ObjectID) ([]*entity.Submission, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"student_id": studentID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var subs []*entity.Submission
	for cursor.Next(ctx) {
		var s entity.Submission
		if err := cursor.Decode(&s); err != nil {
			return nil, err
		}
		subs = append(subs, &s)
	}

	return subs, cursor.Err()
}

func (r *MongoSubmissionRepository) ListSubmissionsByAssignmentAndStudent(ctx context.Context, assignmentID, studentID primitive.ObjectID) ([]*entity.Submission, error) {
	filter := bson.M{
		"assignment_id": assignmentID,
		"student_id":    studentID,
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var subs []*entity.Submission
	for cursor.Next(ctx) {
		var s entity.Submission
		if err := cursor.Decode(&s); err != nil {
			return nil, err
		}
		subs = append(subs, &s)
	}

	return subs, cursor.Err()
}

func (r *MongoSubmissionRepository) ListClassesByStudent(ctx context.Context, studentID primitive.ObjectID) ([]*entity.Class, error) {
	filter := bson.M{"student_ids": studentID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var classes []*entity.Class

	for cursor.Next(ctx) {
		var cl entity.Class
		
		if err := cursor.Decode(&cl); err != nil {
			return nil, err
		}
		classes = append(classes, &cl)
	}

	return classes, cursor.Err()
}