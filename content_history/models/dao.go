package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ContentHistoryModel struct {
	ID            primitive.ObjectID `bson:"_id"`
	ContentId     string             `bson:"content_id"`
	UserId        string             `bson:"user_id"`
	Action        string             `bson:"action"`
	PreviousValue string             `bson:"previous_value"`
	NewValue      string             `bson:"new_value"`
	CreatedAt     time.Time          `bson:"created_at"`
}
