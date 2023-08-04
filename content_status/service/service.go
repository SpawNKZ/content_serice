package service

import (
	"context"
	"github.com/SpawNKZ/content_service/common/errors"
	"github.com/SpawNKZ/content_service/content_status/models"
	repository "github.com/SpawNKZ/content_service/content_status/repo"
)

type Service interface {
	Create(ctx context.Context, data models.CreateRequest) (string, error)
	GetOne(ctx context.Context, data models.IdRequest) (*models.ContentStatus, error)
	GetByName(ctx context.Context, name string) (*models.ContentStatus, error)
	Update(ctx context.Context, data models.UpdateRequest) error
	Delete(ctx context.Context, data models.IdRequest) error
	GetList(ctx context.Context, data models.GetListRequest) ([]*models.ContentStatus, models.Pagination, error)
}

type service struct {
	repo repository.Repository
}

func New(repo repository.Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, data models.CreateRequest) (string, error) {
	return s.repo.Insert(ctx, *data.ToContentStatus())
}

func (s *service) GetOne(ctx context.Context, data models.IdRequest) (*models.ContentStatus, error) {
	return s.repo.FindByID(ctx, data.ID)
}

func (s *service) GetByName(ctx context.Context, name string) (*models.ContentStatus, error) {
	return s.repo.FindByName(ctx, name)
}

func (s *service) GetList(ctx context.Context, data models.GetListRequest) ([]*models.ContentStatus, models.Pagination, error) {
	count, err := s.repo.Count(ctx)

	if err != nil {
		return nil, models.Pagination{}, errors.ErrDB
	}

	contentStatusList, err := s.repo.GetAll(ctx, int64(data.Limit), int64(data.Offset))
	if err != nil {
		return nil, models.Pagination{}, errors.ErrDB
	}

	hasNextPage := int64(data.Offset+data.Limit) < count

	return contentStatusList, models.Pagination{
		Total:      int(count),
		Limit:      data.Limit,
		Offset:     data.Offset,
		IsLastPage: !hasNextPage,
	}, nil
}

func (s *service) Update(ctx context.Context, data models.UpdateRequest) error {
	return s.repo.Update(ctx, data)
}

func (s *service) Delete(ctx context.Context, data models.IdRequest) error {
	return s.repo.DeleteByID(ctx, data.ID)
}
