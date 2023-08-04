package service

import (
	"context"
	"github.com/SpawNKZ/content_service/post/models"
	"github.com/go-kit/log"
	"time"
)

type loggingService struct {
	logger  log.Logger
	service Service
}

func NewLoggingService(logger log.Logger, service Service) Service {
	return &loggingService{
		logger:  logger,
		service: service,
	}
}

func (l *loggingService) Create(ctx context.Context, request models.CreateRequest) (result string, err error) {
	defer func(begin time.Time) {
		l.logger.Log(
			"method", "Create",
			"id", result,
			"err", err,
		)
	}(time.Now())

	return l.service.Create(ctx, request)
}

func (l *loggingService) GetOne(ctx context.Context, request models.IdRequest) (result *models.Post, err error) {
	defer func(begin time.Time) {
		l.logger.Log(
			"method", "GetOne",
			"id", request.ID,
			"err", err,
		)
	}(time.Now())

	return l.service.GetOne(ctx, request)
}

func (l *loggingService) GetList(ctx context.Context, request models.GetListRequest) (result []*models.Post, pagination models.Pagination, err error) {
	defer func(begin time.Time) {
		l.logger.Log(
			"method", "GetList",
			"post", len(result),
			"err", err,
		)
	}(time.Now())

	return l.service.GetList(ctx, request)
}

func (l *loggingService) Update(ctx context.Context, request models.UpdateRequest) (err error) {
	defer func(begin time.Time) {
		l.logger.Log(
			"method", "Update",
			"id", request.ID,
			"err", err,
		)
	}(time.Now())

	return l.service.Update(ctx, request)
}

func (l *loggingService) Delete(ctx context.Context, request models.IdRequest) (err error) {
	defer func(begin time.Time) {
		l.logger.Log(
			"method", "Delete",
			"id", request.ID,
			"err", err,
		)
	}(time.Now())

	return l.service.Delete(ctx, request)
}
