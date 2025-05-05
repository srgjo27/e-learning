package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Course struct {
	ID 				primitive.ObjectID 	 `bson:"_id,omitempty" json:"id"`
	Name 			string             	 `bson:"name" json:"name"`
	Description 	string             	 `bson:"description,omitempty" json:"description,omitempty"`
	AssignedTeacher []primitive.ObjectID `bson:"assigned_teachers" json:"assigned_teachers"`
	CreatedAt 		time.Time		  	 `bson:"created_at" json:"created_at"`
}