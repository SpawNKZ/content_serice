package repository

import (
	"context"
	"github.com/SpawNKZ/content_service/common/errors"
	"github.com/SpawNKZ/content_service/post/models"
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
	Insert(ctx context.Context, postObj models.Post) (string, error)
	FindByID(ctx context.Context, id string) (*models.Post, error)
	Update(ctx context.Context, updateObj models.UpdateRequest) error
	DeleteByID(ctx context.Context, id string) error
	Count(ctx context.Context, filter models.PostFilter) (int64, error)
	GetAll(ctx context.Context, limit int64, offset int64, filter models.PostFilter) ([]*models.Post, error)
}

func NewRepository(client *mongo.Client) Repository {
	return &repository{
		client.Database("core").Collection("post"),
	}
}

func (r *repository) Insert(ctx context.Context, postObj models.Post) (string, error) {
	postModel := postObj.ToRepositoryModel()

	result, err := r.c.InsertOne(ctx, postModel)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), err
}

func (r *repository) FindByID(ctx context.Context, id string) (*models.Post, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{"_id", objectId}, {"deleted_at", nil}}

	var postObj models.PostModel
	err = r.c.FindOne(ctx, filter).Decode(&postObj)

	switch {
	case err == nil:
		var postRes models.Post
		postRes.FromRepositoryModel(&postObj)
		return &postRes, nil
	case err == mongo.ErrNoDocuments:
		return nil, errors.ErrNotFound
	default:
		return nil, err
	}
}

func (r *repository) Count(ctx context.Context, filter models.PostFilter) (int64, error) {
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

func (r *repository) GetAll(ctx context.Context, limit int64, offset int64, filter models.PostFilter) ([]*models.Post, error) {
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

	var postResults []*models.Post

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var post models.PostModel
		if err := cursor.Decode(&post); err != nil {
			return nil, err
		} else {
			var postRes models.Post
			postRes.FromRepositoryModel(&post)
			postResults = append(postResults, &postRes)
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return postResults, nil
}

func (r *repository) Update(ctx context.Context, updateObj models.UpdateRequest) error {
	objectId, err := primitive.ObjectIDFromHex(updateObj.ID)
	if err != nil {
		return err
	}
	filter := bson.D{{"_id", objectId}}

	updateMap := updateObj.ToPost()

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

	updateResult, err := r.c.UpdateOne(ctx, filter, bson.D{{"$set", bson.D{{"deleted_at", time.Now()}}}})
	if err != nil {
		return err
	}
	if updateResult.ModifiedCount == 0 {
		return errors.ErrNotFound
	}
	return nil
}
