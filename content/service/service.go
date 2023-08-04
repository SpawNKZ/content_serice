package service

import (
	"context"
	"github.com/SpawNKZ/content_service/common/errors"
	"github.com/SpawNKZ/content_service/content/models"
	repository "github.com/SpawNKZ/content_service/content/repo"
	contentHistoryModel "github.com/SpawNKZ/content_service/content_history/models"
	contentHistory "github.com/SpawNKZ/content_service/content_history/service"
	contentStatusModel "github.com/SpawNKZ/content_service/content_status/models"
	contentStatus "github.com/SpawNKZ/content_service/content_status/service"
	"time"
)

const (
	actionCreate       = "create"
	actionUpdate       = "update"
	actionAssignAuthor = "assign_author"
	actionDelete       = "delete"
	actionChangeStatus = "change_status"
)

const (
	statusDraft     = "draft"
	statusArchived  = "archived"
	statusPublished = "published"
)

type Service interface {
	Create(ctx context.Context, data models.CreateRequest) (string, error)
	GetOne(ctx context.Context, data models.IdRequest) (*models.Content, error)
	Update(ctx context.Context, data models.UpdateRequest) error
	AssignAuthor(ctx context.Context, data models.AssignAuthorRequest) error
	ChangeStatus(ctx context.Context, data models.ChangeStatusRequest) error
	Delete(ctx context.Context, data models.IdRequest) error
	GetList(ctx context.Context, data models.GetListRequest) ([]*models.Content, models.Pagination, error)
}

type service struct {
	repo       repository.Repository
	history    contentHistory.Service
	subject    repository.SubjectRepository
	microtopic repository.MicrotopicRepository
	status     contentStatus.Service
}

func New(repo repository.Repository, subject repository.SubjectRepository, microtopic repository.MicrotopicRepository, history contentHistory.Service, status contentStatus.Service) Service {
	return &service{
		repo:       repo,
		subject:    subject,
		history:    history,
		microtopic: microtopic,
		status:     status,
	}
}

func (s *service) Create(ctx context.Context, data models.CreateRequest) (string, error) {
	_, err := s.subject.GetSubjectId(ctx, data.SubjectId)
	if err != nil {
		return "", err
	}

	_, err = s.microtopic.GetMicrotopicId(ctx, data.MicrotopicId)
	if err != nil {
		return "", err
	}

	contentStatusObj, err := s.status.GetByName(ctx, statusDraft)
	if err != nil {
		return "", err
	}

	data.StatusId = contentStatusObj.Name
	contentId, err := s.repo.Insert(ctx, *data.ToContent())
	if err != nil {
		return "", err
	}

	contentObj, err := s.GetOne(ctx, models.IdRequest{ID: contentId})
	if err != nil {
		return "", err
	}

	err = s.history.Create(ctx, contentHistoryModel.ContentHistory{
		ContentId: contentId,
		UserId:    contentObj.AuthorId,
		Action:    actionCreate,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return "", err
	}
	return contentId, nil
}

func (s *service) GetOne(ctx context.Context, data models.IdRequest) (*models.Content, error) {
	contentObj, err := s.repo.FindByID(ctx, data.ID)
	if err != nil {
		return nil, err
	}

	subject, err := s.subject.GetSubjectId(ctx, contentObj.SubjectId)
	if err != nil {
		return nil, err
	}

	var contentRes models.Content
	contentRes.FromRepositoryModel(contentObj)
	contentRes.Subject = subject

	return &contentRes, nil
}

func (s *service) GetList(ctx context.Context, data models.GetListRequest) ([]*models.Content, models.Pagination, error) {
	count, err := s.repo.Count(ctx, data.ContentFilter)

	if err != nil {
		return nil, models.Pagination{}, errors.ErrDB
	}

	contentList, err := s.repo.GetAll(ctx, int64(data.Limit), int64(data.Offset), data.ContentFilter)
	if err != nil {
		return nil, models.Pagination{}, errors.ErrDB
	}

	hasNextPage := int64(data.Offset+data.Limit) < count

	return contentList, models.Pagination{
		Total:      int(count),
		Limit:      data.Limit,
		Offset:     data.Offset,
		IsLastPage: !hasNextPage,
	}, nil
}

func (s *service) Update(ctx context.Context, data models.UpdateRequest) error {
	contentObj, err := s.GetOne(ctx, models.IdRequest{ID: data.ID})
	if err != nil {
		return err
	}

	err = s.repo.Update(ctx, data)
	if err != nil {
		return err
	}

	err = s.history.Create(ctx, contentHistoryModel.ContentHistory{
		ContentId: data.ID,
		UserId:    contentObj.AuthorId,
		Action:    actionUpdate,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *service) AssignAuthor(ctx context.Context, data models.AssignAuthorRequest) error {
	contentObj, err := s.GetOne(ctx, models.IdRequest{ID: data.ID})
	if err != nil {
		return err
	}

	err = s.repo.UpdateAuthor(ctx, data)
	if err != nil {
		return err
	}

	err = s.history.Create(ctx, contentHistoryModel.ContentHistory{
		ContentId: data.ID,
		UserId:    contentObj.AuthorId,
		Action:    actionAssignAuthor,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *service) ChangeStatus(ctx context.Context, data models.ChangeStatusRequest) error {
	contentObj, err := s.GetOne(ctx, models.IdRequest{ID: data.ID})
	if err != nil {
		return err
	}

	_, err = s.status.GetOne(ctx, contentStatusModel.IdRequest{ID: data.StatusId})
	if err != nil {
		return err
	}

	err = s.repo.UpdateStatus(ctx, data)
	if err != nil {
		return err
	}

	err = s.history.Create(ctx, contentHistoryModel.ContentHistory{
		ContentId: data.ID,
		UserId:    contentObj.AuthorId,
		Action:    actionChangeStatus,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *service) Delete(ctx context.Context, data models.IdRequest) error {
	contentObj, err := s.GetOne(ctx, models.IdRequest{ID: data.ID})
	if err != nil {
		return err
	}

	statusObj, err := s.status.GetOne(ctx, contentStatusModel.IdRequest{ID: data.ID})
	if err != nil {
		return err
	}

	if statusObj.IsRemovable == false {
		return errors.ErrContentIsNotRemovable
	}

	//TODO: content status should be archived

	err = s.repo.DeleteByID(ctx, data.ID)
	if err != nil {
		return err
	}

	err = s.history.Create(ctx, contentHistoryModel.ContentHistory{
		ContentId: data.ID,
		UserId:    contentObj.AuthorId,
		Action:    actionDelete,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}
