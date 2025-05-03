package entity

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string

const (
	RoleAdmin 	Role = "admin"
	RoleTeacher Role = "teacher"
	RoleStudent Role = "student"
)

type User struct {
	ID 		  primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email 	  string 		   	 `bson:"email" json:"email"`
	Password  string			 `bson:"password" json:"password"`
	Role 	  Role			 	 `bson:"role" json:"role"`
	CreatedAt time.Time			 `bson:"created_at" json:"created_at"` 			
}

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrEmailExists     = errors.New("email already exists")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidToken	   = errors.New("invalid token")
	ErrUnauthorized	   = errors.New("unauthorized access")
)