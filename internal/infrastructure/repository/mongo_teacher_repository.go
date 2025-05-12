package repository

import (
	"context"

	"github.com/srgjo27/e-learning/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *MongoCourseRepository) ListCoursesByTeacher(ctx context.Context, teacherID primitive.ObjectID) ([]*entity.Course, error) {
	filter := bson.M{"assigned_teachers": teacherID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var courses []*entity.Course
	for cursor.Next(ctx) {
		var c entity.Course

		if err := cursor.Decode(&c); err != nil {
			return nil, err
		}

		courses = append(courses, &c)
	}

	return courses, cursor.Err()
}

func (r *MongoClassRepository) ListClassesByTeacher(ctx context.Context, teacherID primitive.ObjectID) ([]*entity.Class, error) {
	filter := bson.M{"teacher_ids": teacherID}
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

func (r *MongoUserRepository) FindUsersByIDs(ctx context.Context, ids []primitive.ObjectID) ([]*entity.User, error) {
	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var users []*entity.User
	for cursor.Next(ctx) {
		var u entity.User
		if err := cursor.Decode(&u); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, cursor.Err()
}