package repository

import (
	"context"
	"encoding/json"
	"github.com/SpawNKZ/content_service/common/errors"
	"github.com/SpawNKZ/content_service/content/models"
	natsTransport "github.com/go-kit/kit/transport/nats"
	"github.com/nats-io/nats.go"
)

type microtopicRepository struct {
	nc *nats.Conn
}

type MicrotopicRepository interface {
	GetMicrotopicId(ctx context.Context, subjectId int64) (*models.Microtopic, error)
}

func NewMicrotopicRepository(nc *nats.Conn) MicrotopicRepository {
	return &microtopicRepository{nc: nc}
}

type getMicrotopicByIdRequest struct {
	ID int64 `json:"id"`
}

type getMicrotopicByIdResponse struct {
	Microtopic *models.Microtopic `json:"microtopic,omitempty"`
	Err        error              `json:"error,omitempty"`
}

func (s *microtopicRepository) GetMicrotopicId(ctx context.Context, microtopicId int64) (*models.Microtopic, error) {
	publisher := natsTransport.NewPublisher(
		s.nc,
		"microtopics.GetById",
		natsTransport.EncodeJSONRequest,
		decodeMicrotopicById,
	)

	res, err := publisher.Endpoint()(ctx, getMicrotopicByIdRequest{ID: microtopicId})
	if err != nil {
		return nil, err
	}

	response, ok := res.(getMicrotopicByIdResponse)
	if !ok {
		return nil, errors.ErrNoResponseFromSubject
	}

	if response.Microtopic == nil {
		return nil, errors.ErrNoResponseFromSubject
	}
	return response.Microtopic, nil

}

func decodeMicrotopicById(_ context.Context, msg *nats.Msg) (response interface{}, err error) {
	var res getMicrotopicByIdResponse

	err = json.Unmarshal(msg.Data, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
