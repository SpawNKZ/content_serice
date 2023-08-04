package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ContentModel struct {
	ID           primitive.ObjectID `bson:"_id"`
	Locale       string             `bson:"locale"`
	Body         string             `bson:"body"`
	Description  string             `bson:"description"`
	Resources    []string           `bson:"resources"`
	SubjectId    int64              `bson:"subject_id"`
	MicrotopicId int64              `bson:"microtopic_id"`
	StatusId     string             `bson:"status_id"`
	AuthorId     string             `bson:"author_id"`
	Difficulty   int                `bson:"difficulty"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
	DeletedAt    *time.Time         `bson:"deleted_at"`
}
