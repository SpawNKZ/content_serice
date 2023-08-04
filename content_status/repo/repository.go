package repository

import (
	"context"
	"github.com/SpawNKZ/content_service/common/errors"
	"github.com/SpawNKZ/content_service/content_status/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type repository struct {
	c *mongo.Collection
}

type Repository interface {
	Insert(ctx context.Context, contentStatusObj models.ContentStatus) (string, error)
	FindByID(ctx context.Context, id string) (*models.ContentStatus, error)
	FindByName(ctx context.Context, name string) (*models.ContentStatus, error)
	Update(ctx context.Context, updateObj models.UpdateRequest) error
	DeleteByID(ctx context.Context, id string) error
	Count(ctx context.Context) (int64, error)
	GetAll(ctx context.Context, limit int64, offset int64) ([]*models.ContentStatus, error)
}

func NewRepository(client *mongo.Client) Repository {
	return &repository{
		client.Database("core").Collection("content_status"),
	}
}

func (r *repository) Insert(ctx context.Context, contentStatusObj models.ContentStatus) (string, error) {
	contentStatusModel := contentStatusObj.ToRepositoryModel()

	result, err := r.c.InsertOne(ctx, contentStatusModel)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), err
}

func (r *repository) FindByID(ctx context.Context, id string) (*models.ContentStatus, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{"_id", objectId}}

	var contentStatusObj models.ContentStatusModel
	err = r.c.FindOne(ctx, filter).Decode(&contentStatusObj)

	switch {
	case err == nil:
		var contentStatusRes models.ContentStatus
		contentStatusRes.FromRepositoryModel(&contentStatusObj)
		return &contentStatusRes, nil
	case err == mongo.ErrNoDocuments:
		return nil, errors.ErrNotFound
	default:
		return nil, err
	}
}

func (r *repository) FindByName(ctx context.Context, name string) (*models.ContentStatus, error) {
	filter := bson.D{{"name", name}}

	var contentStatusObj models.ContentStatusModel
	err := r.c.FindOne(ctx, filter).Decode(&contentStatusObj)

	switch {
	case err == nil:
		var contentStatusRes models.ContentStatus
		contentStatusRes.FromRepositoryModel(&contentStatusObj)
		return &contentStatusRes, nil
	case err == mongo.ErrNoDocuments:
		return nil, errors.ErrNotFound
	default:
		return nil, err
	}
}

func (r *repository) Count(ctx context.Context) (int64, error) {
	filterQuery := bson.M{}
	count, err := r.c.CountDocuments(ctx, filterQuery)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *repository) GetAll(ctx context.Context, limit int64, offset int64) ([]*models.ContentStatus, error) {
	filterQuery := bson.M{}

	opts := options.Find()
	opts.SetSkip(offset)
	opts.SetLimit(limit)

	cursor, err := r.c.Find(ctx, filterQuery, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var contentStatusResults []*models.ContentStatus

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var contentStatus models.ContentStatusModel
		if err := cursor.Decode(&contentStatus); err != nil {
			return nil, err
		} else {
			var contentRes models.ContentStatus
			contentRes.FromRepositoryModel(&contentStatus)
			contentStatusResults = append(contentStatusResults, &contentRes)
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return contentStatusResults, nil
}

func (r *repository) Update(ctx context.Context, updateObj models.UpdateRequest) error {
	objectId, err := primitive.ObjectIDFromHex(updateObj.ID)
	if err != nil {
		return err
	}
	filter := bson.D{{"_id", objectId}}

	updateMap := updateObj.ToContentStatus()

	var update bson.D
	for key, value := range updateMap {
		update = append(update, bson.E{Key: key, Value: value})
	}

	updateResult, err := r.c.UpdateOne(ctx, filter, bson.D{{"$set", update}})
	if err != nil {
		return err
	}
	if updateResult.ModifiedCount == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *repository) DeleteByID(ctx context.Context, id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.D{{"_id", objectId}}

	result, err := r.c.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.ErrNotFound
	}
	return nil
}
