package repository

import (
	"context"
	"github.com/SpawNKZ/content_service/common/errors"
	"github.com/SpawNKZ/content_service/content/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"time"
)

type repository struct {
	c *mongo.Collection
}

type Repository interface {
	Insert(ctx context.Context, contentObj models.Content) (string, error)
	FindByID(ctx context.Context, id string) (*models.ContentModel, error)
	UpdateStatus(ctx context.Context, updateObj models.ChangeStatusRequest) error
	Update(ctx context.Context, updateObj models.UpdateRequest) error
	UpdateAuthor(ctx context.Context, updateObj models.AssignAuthorRequest) error
	DeleteByID(ctx context.Context, id string) error
	Count(ctx context.Context, filter models.ContentFilter) (int64, error)
	GetAll(ctx context.Context, limit int64, offset int64, filter models.ContentFilter) ([]*models.Content, error)
}

func NewRepository(client *mongo.Client) Repository {
	return &repository{
		client.Database("core").Collection("content"),
	}
}

func (r *repository) Insert(ctx context.Context, contentObj models.Content) (string, error) {
	contentModel := contentObj.ToRepositoryModel()

	result, err := r.c.InsertOne(ctx, contentModel)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), err
}

func (r *repository) FindByID(ctx context.Context, id string) (*models.ContentModel, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{"_id", objectId}, {"deleted_at", nil}}

	var contentObj models.ContentModel
	err = r.c.FindOne(ctx, filter).Decode(&contentObj)

	switch {
	case err == nil:
		return &contentObj, nil
	case err == mongo.ErrNoDocuments:
		return nil, errors.ErrNotFound
	default:
		return nil, err
	}
}

func (r *repository) Count(ctx context.Context, filter models.ContentFilter) (int64, error) {
	filterQuery := bson.M{"deleted_at": nil}
	filterValue := reflect.ValueOf(filter)
	filterType := reflect.TypeOf(filter)

	for i := 0; i < filterValue.NumField(); i++ {
		field := filterValue.Field(i)
		if !field.IsZero() {
			fieldName := filterType.Field(i).Tag.Get("bson")
			filterQuery[fieldName] = field.Interface()
		}
	}

	count, err := r.c.CountDocuments(ctx, filterQuery)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *repository) GetAll(ctx context.Context, limit int64, offset int64, filter models.ContentFilter) ([]*models.Content, error) {
	filterQuery := bson.M{"deleted_at": nil}
	filterValue := reflect.ValueOf(filter)
	filterType := reflect.TypeOf(filter)

	for i := 0; i < filterValue.NumField(); i++ {
		field := filterValue.Field(i)
		if !field.IsZero() {
			fieldName := filterType.Field(i).Tag.Get("bson")
			filterQuery[fieldName] = field.Interface()
		}
	}

	options := options.Find()
	options.SetSkip(offset)
	options.SetLimit(limit)

	cursor, err := r.c.Find(ctx, filterQuery, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var contentResults []*models.Content

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var content models.ContentModel
		if err := cursor.Decode(&content); err != nil {
			return nil, err
		} else {
			var contentRes models.Content
			contentRes.FromRepositoryModel(&content)
			contentResults = append(contentResults, &contentRes)
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return contentResults, nil
}

func (r *repository) Update(ctx context.Context, updateObj models.UpdateRequest) error {
	objectId, err := primitive.ObjectIDFromHex(updateObj.ID)
	if err != nil {
		return err
	}
	filter := bson.D{{"_id", objectId}}

	updateMap := updateObj.ToContent()

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

func (r *repository) UpdateAuthor(ctx context.Context, updateObj models.AssignAuthorRequest) error {
	objectId, err := primitive.ObjectIDFromHex(updateObj.ID)
	if err != nil {
		return err
	}
	filter := bson.D{{"_id", objectId}}

	updateMap := updateObj.ToContent()

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

	updateResult, err := r.c.UpdateOne(ctx, filter, bson.D{{"$set", bson.M{"deleted_at": time.Now()}}})
	if err != nil {
		return err
	}
	if updateResult.ModifiedCount == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *repository) UpdateStatus(ctx context.Context, updateObj models.ChangeStatusRequest) error {
	objectId, err := primitive.ObjectIDFromHex(updateObj.ID)
	if err != nil {
		return err
	}
	filter := bson.D{{"_id", objectId}}

	updateMap := updateObj.ToContent()

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
