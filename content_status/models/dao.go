package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ContentStatusModel struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `bson:"name"`
	IsRemovable bool               `bson:"is_removable"`
}
