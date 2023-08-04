package endpoints

import (
	"context"
	"github.com/SpawNKZ/content_service/content/models"
	contentService "github.com/SpawNKZ/content_service/content/service"
	"github.com/go-kit/kit/endpoint"
)

func MakeCreateEndpoint(s contentService.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(models.CreateRequest)
		id, err := s.Create(ctx, req)
		return models.CreateResponse{ID: id, Err: err}, nil
	}
}

func MakeGetOneEndpoint(s contentService.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(models.IdRequest)
		res, e := s.GetOne(ctx, req)
		return models.GetOneResponse{Content: res, Err: e}, nil
	}
}

func MakeUpdateEndpoint(s contentService.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(models.UpdateRequest)
		e := s.Update(ctx, req)
		return models.UpdateResponse{Err: e}, nil
	}
}

func MakeAssignAuthorEndpoint(s contentService.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(models.AssignAuthorRequest)
		e := s.AssignAuthor(ctx, req)
		return models.UpdateResponse{Err: e}, nil
	}
}

func MakeChangeStatusEndpoint(s contentService.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(models.ChangeStatusRequest)
		e := s.ChangeStatus(ctx, req)
		return models.UpdateResponse{Err: e}, nil
	}
}

func MakeDeleteOneEndpoint(s contentService.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(models.IdRequest)
		e := s.Delete(ctx, req)
		return models.DeleteOneResponse{Err: e}, nil
	}
}

func MakeGetListEndpoint(s contentService.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(models.GetListRequest)
		res, pagination, e := s.GetList(ctx, req)
		return models.GetListResponse{ContentList: res, Pagination: pagination, Err: e}, nil
	}
}
