package repository

import (
	"context"
	"errors"

	"github.com/srgjo27/e-learning/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoCourseRepository struct {
	collection *mongo.Collection
}

type MongoClassRepository struct {
	collection *mongo.Collection
}

type MongoAnnouncementRepository struct {
	collection *mongo.Collection
}

func NewMongoCourseRepository(c *mongo.Collection) *MongoCourseRepository {
	return &MongoCourseRepository{collection: c}
}

func NewMongoClassRepository(c *mongo.Collection) *MongoClassRepository {
	return &MongoClassRepository{collection: c}
}

func NewMongoAnnouncementRepository(c *mongo.Collection) *MongoAnnouncementRepository {
	return &MongoAnnouncementRepository{collection: c}
}

// --- Course ---
func (r *MongoCourseRepository) CreateCourse(ctx context.Context, course *entity.Course) error {
	course.CreatedAt = course.CreatedAt.UTC()
	_, err := r.collection.InsertOne(ctx, course)
	return err
}

func (r *MongoCourseRepository) GetCourse(ctx context.Context, id string) (*entity.Course, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var course entity.Course

	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&course)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, entity.ErrCourseNotFound
		}

		return nil, err
	}

	return &course, nil
}

func (r *MongoCourseRepository) UpdateCourse(ctx context.Context, course *entity.Course) error {
	if course.ID.IsZero() {
		return errors.New("course id required")
	}

	filter := bson.M{"_id": course.ID}
	update := bson.M{
		"$set": bson.M{
			"name": 		   	 course.Name,
			"description":    	 course.Description,
			"assigned_teachers": course.AssignedTeacher,
		},
	}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.New("course not found")
	}

	return nil
}

func (r *MongoCourseRepository) DeleteCourse(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("course not found")
	}

	return nil
}

func (r *MongoCourseRepository) ListCourses(ctx context.Context) ([]*entity.Course, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var courses []*entity.Course

	for cursor.Next(ctx) {
		var course entity.Course

		if err := cursor.Decode(&course); err != nil {
			return nil, err
		}
		courses = append(courses, &course)
	}

	return courses, cursor.Err()
}

// --- Class ---
func (r *MongoClassRepository) CreateClass(ctx context.Context, class *entity.Class) error {
	class.CreatedAt = class.CreatedAt.UTC()
	_, err := r.collection.InsertOne(ctx, class)
	return err
}

func (r *MongoClassRepository) GetClass(ctx context.Context, id string) (*entity.Class, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var class entity.Class

	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&class)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("class not found")
		}

		return nil, err
	}

	return &class, nil
}

func (r *MongoClassRepository) UpdateClass(ctx context.Context, class *entity.Class) error {
	if class.ID.IsZero() {
		return errors.New("class id required")
	}

	filter := bson.M{"_id": class.ID}
	update := bson.M{
		"$set": bson.M{
			"name":        class.Name,
			"student_ids": class.StudentIDs,
			"teacher_ids": class.TeacherIDs,
		},
	}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.New("class not found")
	}

	return nil
}

func (r *MongoClassRepository) DeleteClass(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("class not found")
	}

	return nil
}

func (r *MongoClassRepository) ListClasses(ctx context.Context) ([]*entity.Class, error) {
	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var classes []*entity.Class
	for cursor.Next(ctx) {
		var class entity.Class
		if err := cursor.Decode(&class); err != nil {
			return nil, err
		}
		classes = append(classes, &class)
	}
	return classes, cursor.Err()
}

// --- Announcement ---
func (r *MongoAnnouncementRepository) CreateAnnouncement(ctx context.Context, ann *entity.Announcement) error {
	ann.CreatedAt = ann.CreatedAt.UTC()
	_, err := r.collection.InsertOne(ctx, ann)
	return err
}

func (r *MongoAnnouncementRepository) GetAnnouncement(ctx context.Context, id string) (*entity.Announcement, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var ann entity.Announcement
	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&ann)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("announcement not found")
		}
		return nil, err
	}
	return &ann, nil
}

func (r *MongoAnnouncementRepository) UpdateAnnouncement(ctx context.Context, ann *entity.Announcement) error {
	if ann.ID.IsZero() {
		return errors.New("announcement id required")
	}
	filter := bson.M{"_id": ann.ID}
	update := bson.M{
		"$set": bson.M{
			"title":           ann.Title,
			"content":         ann.Content,
			"target_audience": ann.TargetAudience,
			"target_id":       ann.TargetID,
		},
	}
	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("announcement not found")
	}
	return nil
}

func (r *MongoAnnouncementRepository) DeleteAnnouncement(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("announcement not found")
	}
	return nil
}

func (r *MongoAnnouncementRepository) ListAnnouncements(ctx context.Context) ([]*entity.Announcement, error) {
	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var anns []*entity.Announcement
	for cursor.Next(ctx) {
		var ann entity.Announcement
		if err := cursor.Decode(&ann); err != nil {
			return nil, err
		}
		anns = append(anns, &ann)
	}
	return anns, cursor.Err()
}