package repository

import (
	"context"
	"encoding/json"
	"github.com/SpawNKZ/content_service/common/errors"
	"github.com/SpawNKZ/content_service/content/models"
	natsTransport "github.com/go-kit/kit/transport/nats"
	"github.com/nats-io/nats.go"
)

type subjectRepository struct {
	nc *nats.Conn
}

type SubjectRepository interface {
	GetSubjectId(ctx context.Context, subjectId int64) (*models.Subject, error)
}

func NewSubjectRepository(nc *nats.Conn) SubjectRepository {
	return &subjectRepository{nc: nc}
}

type getSubjectByIdRequest struct {
	ID int64 `json:"id"`
}
type getSubjectByIdResponse struct {
	Subject *models.Subject `json:"subject,omitempty"`
	Err     error           `json:"error,omitempty"`
}

func (s *subjectRepository) GetSubjectId(ctx context.Context, subjectId int64) (*models.Subject, error) {
	publisher := natsTransport.NewPublisher(
		s.nc,
		"subjects.GetById",
		natsTransport.EncodeJSONRequest,
		decodeSubjectById,
	)

	res, err := publisher.Endpoint()(ctx, getSubjectByIdRequest{ID: subjectId})
	if err != nil {
		return nil, err
	}

	response, ok := res.(getSubjectByIdResponse)
	if !ok {
		return nil, errors.ErrNoResponseFromSubject
	}

	if response.Subject == nil {
		return nil, errors.ErrNoResponseFromSubject
	}
	return response.Subject, nil

}

func decodeSubjectById(_ context.Context, msg *nats.Msg) (response interface{}, err error) {
	var res getSubjectByIdResponse

	err = json.Unmarshal(msg.Data, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
