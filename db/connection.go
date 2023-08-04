package db

import (
	"context"
	"github.com/go-kit/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDBConnection(ctx context.Context, dbSource string, logger log.Logger) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbSource))
	if err != nil {
		logger.Log("cannot connect to DB", err)
		return nil, err
	}

	return client, nil
}
