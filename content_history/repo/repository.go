package repository

import (
	"context"
	"github.com/SpawNKZ/content_service/content_history/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type repository struct {
	c *mongo.Collection
}

type Repository interface {
	Insert(ctx context.Context, contentStatusObj models.ContentHistory) error
}

func NewRepository(client *mongo.Client) Repository {
	return &repository{
		client.Database("core").Collection("content_history"),
	}
}

func (r *repository) Insert(ctx context.Context, contentStatusObj models.ContentHistory) error {
	contentStatusModel := contentStatusObj.ToRepositoryModel()

	_, err := r.c.InsertOne(ctx, contentStatusModel)
	if err != nil {
		return err
	}

	return nil
}
