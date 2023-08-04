package service

import (
	"context"
	"github.com/SpawNKZ/content_service/content_history/models"
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

func (l *loggingService) Create(ctx context.Context, request models.ContentHistory) (err error) {
	defer func(begin time.Time) {
		l.logger.Log(
			"method", "Create",
			"err", err,
		)
	}(time.Now())

	return l.service.Create(ctx, request)
}
