package entity

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Submission struct {
	ID 				primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AssignmentID 	primitive.ObjectID `bson:"assignment_id" json:"assignment_id"`
	StudentID 		primitive.ObjectID `bson:"student_id" json:"student_id"`
	Content			string			   `bson:"content" json:"content"` // text submission or link
	SubmittedAt 	time.Time 		   `bson:"submitted_at" json:"submitted_at"`
	Grade			*float64 		   `bson:"grade,omitempty" json:"grade,omitempty"`
	Feedback		*string			   `bson:"feedback,omitempty" json:"feedback,omitempty"`
}

var (
	ErrSubmissionNotFound = errors.New("submission not found")
)