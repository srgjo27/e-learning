package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Assignment struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	CourseID    primitive.ObjectID `bson:"course_id" json:"course_id"`
	DueDate     time.Time          `bson:"due_date" json:"due_date"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}