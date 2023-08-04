package models

import (
	"github.com/go-kit/kit/endpoint"
)

var (
	_ endpoint.Failer = CreateResponse{}
)

type CreateResponse struct {
	ID  string `json:"id"`
	Err error  `json:"error,omitempty"`
}

func (r CreateResponse) Failed() error { return r.Err }
