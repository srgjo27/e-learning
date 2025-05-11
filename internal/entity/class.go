package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Class struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name       string               `bson:"name" json:"name"`
	StudentIDs []primitive.ObjectID `bson:"student_ids" json:"student_ids"`
	TeacherIDs []primitive.ObjectID `bson:"teacher_ids" json:"teacher_ids"`
	CourseID   primitive.ObjectID   `bson:"course_id" json:"course_id"`
	CreatedAt  time.Time            `bson:"created_at" json:"created_at"`
}