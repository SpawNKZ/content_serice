package service

import (
	"context"
	"github.com/SpawNKZ/content_service/common/errors"
	contentModel "github.com/SpawNKZ/content_service/content/models"
	contentRepository "github.com/SpawNKZ/content_service/content/service"
	"github.com/SpawNKZ/content_service/post/models"
	repository "github.com/SpawNKZ/content_service/post/repo"
)

type Service interface {
	Create(ctx context.Context, data models.CreateRequest) (string, error)
	GetOne(ctx context.Context, data models.IdRequest) (*models.Post, error)
	Update(ctx context.Context, data models.UpdateRequest) error
	Delete(ctx context.Context, data models.IdRequest) error
	GetList(ctx context.Context, data models.GetListRequest) ([]*models.Post, models.Pagination, error)
}

type service struct {
	repo       repository.Repository
	contentSvc contentRepository.Service
}

func New(repo repository.Repository, contentSvc contentRepository.Service) Service {
	return &service{
		repo:       repo,
		contentSvc: contentSvc,
	}
}

func (s *service) Create(ctx context.Context, data models.CreateRequest) (string, error) {
	_, err := s.contentSvc.GetOne(ctx, contentModel.IdRequest{ID: data.ContentId})
	if err != nil {
		return "", err
	}

	return s.repo.Insert(ctx, *data.ToPost())
}

func (s *service) GetOne(ctx context.Context, data models.IdRequest) (*models.Post, error) {
	return s.repo.FindByID(ctx, data.ID)
}

func (s *service) GetList(ctx context.Context, data models.GetListRequest) ([]*models.Post, models.Pagination, error) {
	count, err := s.repo.Count(ctx, data.PostFilter)

	if err != nil {
		return nil, models.Pagination{}, errors.ErrDB
	}

	postList, err := s.repo.GetAll(ctx, int64(data.Limit), int64(data.Offset), data.PostFilter)
	if err != nil {
		return nil, models.Pagination{}, errors.ErrDB
	}

	hasNextPage := int64(data.Offset+data.Limit) < count

	return postList, models.Pagination{
		Total:      int(count),
		Limit:      data.Limit,
		Offset:     data.Offset,
		IsLastPage: !hasNextPage,
	}, nil
}

func (s *service) Update(ctx context.Context, data models.UpdateRequest) error {
	_, err := s.GetOne(ctx, models.IdRequest{ID: data.ID})
	if err != nil {
		return err
	}
	return s.repo.Update(ctx, data)
}

func (s *service) Delete(ctx context.Context, data models.IdRequest) error {
	_, err := s.GetOne(ctx, models.IdRequest{ID: data.ID})
	if err != nil {
		return err
	}
	return s.repo.DeleteByID(ctx, data.ID)
}
