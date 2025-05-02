package entity

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID 		  primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email 	  string 		   	 `bson:"email" json:"email"`
	Password  string			 `bson:"password" json:"password"`
	CreatedAt time.Time			 `bson:"created_at" json:"created_at"` 			
}

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrEmailExists     = errors.New("email already exists")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidToken	   = errors.New("invalid token")
)