package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TargetAudience string

const (
	AudienceAll    TargetAudience = "all"
	AudienceClass  TargetAudience = "class"
	AudienceCourse TargetAudience = "course"
)

type Announcement struct {
	ID 				primitive.ObjectID 	`bson:"_id,omitempty" json:"id"`
	Title 			string 				`bson:"title" json:"title"`
	Content 		string 				`bson:"content" json:"content"`
	TargetAudience 	TargetAudience 		`bson:"target_audience" json:"target_audience"`
	TargetID 		primitive.ObjectID  `bson:"target_id,omitempty" json:"target_id,omitempty"` // e.g. class or course ID when targeted
	CreatedAt 		time.Time 			`bson:"created_at" json:"created_at"`
}