package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PostModel struct {
	ID          primitive.ObjectID `bson:"_id"`
	Category    string             `bson:"category"`
	Resources   []string           `bson:"resources"`
	ContentId   string             `bson:"content_id"`
	Description string             `bson:"description"`
	DeletedAt   *time.Time         `bson:"deleted_at"`
}
