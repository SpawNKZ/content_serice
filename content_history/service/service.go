package service

import (
	"context"
	"github.com/SpawNKZ/content_service/content_history/models"
	repository "github.com/SpawNKZ/content_service/content_history/repo"
)

type Service interface {
	Create(ctx context.Context, data models.ContentHistory) error
}

type service struct {
	repo repository.Repository
}

func New(repo repository.Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, data models.ContentHistory) error {
	return s.repo.Insert(ctx, data)
}
