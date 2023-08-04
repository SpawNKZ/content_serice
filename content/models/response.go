package models

import (
	"github.com/go-kit/kit/endpoint"
)

type Pagination struct {
	Total      int  `json:"total"`
	Limit      int  `json:"limit"`
	Offset     int  `json:"offset"`
	IsLastPage bool `json:"isLastPage"`
}

const (
	DefaultLimit = 10
	MaxOffset    = 2147483647
	MaxLimit     = 1000
)

var (
	_ endpoint.Failer = CreateResponse{}
	_ endpoint.Failer = GetOneResponse{}
	_ endpoint.Failer = UpdateResponse{}
	_ endpoint.Failer = DeleteOneResponse{}
	_ endpoint.Failer = GetListResponse{}
)

type CreateResponse struct {
	ID  string `json:"id"`
	Err error  `json:"error,omitempty"`
}

func (r CreateResponse) Failed() error { return r.Err }

type GetOneResponse struct {
	Content *Content `json:"content,omitempty"`
	Err     error    `json:"error,omitempty"`
}

func (r GetOneResponse) Failed() error { return r.Err }

type UpdateResponse struct {
	Err error `json:"error,omitempty"`
}

func (r UpdateResponse) Failed() error { return r.Err }

type DeleteOneResponse struct {
	Err error `json:"error,omitempty"`
}

func (r DeleteOneResponse) Failed() error { return r.Err }

type GetListResponse struct {
	ContentList []*Content `json:"content,omitempty"`
	Pagination
	Err error `json:"error,omitempty"`
}

func (r GetListResponse) Failed() error { return r.Err }
