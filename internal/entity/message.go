package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	SenderID    primitive.ObjectID   `bson:"sender_id" json:"sender_id"`
	ReceiverIDs []primitive.ObjectID `bson:"receiver_ids" json:"receiver_ids"` // Recipients: students or classes
	Content     string               `bson:"content" json:"content"`
	CreatedAt   time.Time            `bson:"created_at" json:"created_at"`
}