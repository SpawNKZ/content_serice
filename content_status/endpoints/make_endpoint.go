package endpoints

import (
	"context"
	"github.com/SpawNKZ/content_service/content_status/models"
	contentStatusService "github.com/SpawNKZ/content_service/content_status/service"
	"github.com/go-kit/kit/endpoint"
)

func MakeCreateEndpoint(s contentStatusService.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(models.CreateRequest)
		id, err := s.Create(ctx, req)
		return models.CreateResponse{ID: id, Err: err}, nil
	}
}

func MakeGetOneEndpoint(s contentStatusService.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(models.IdRequest)
		res, e := s.GetOne(ctx, req)
		return models.GetOneResponse{ContentStatus: res, Err: e}, nil
	}
}

func MakeUpdateEndpoint(s contentStatusService.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(models.UpdateRequest)
		e := s.Update(ctx, req)
		return models.UpdateResponse{Err: e}, nil
	}
}

func MakeDeleteOneEndpoint(s contentStatusService.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(models.IdRequest)
		e := s.Delete(ctx, req)
		return models.DeleteOneResponse{Err: e}, nil
	}
}

func MakeGetListEndpoint(s contentStatusService.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(models.GetListRequest)
		res, pagination, e := s.GetList(ctx, req)
		return models.GetListResponse{ContentStatusList: res, Pagination: pagination, Err: e}, nil
	}
}
