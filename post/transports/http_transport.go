package transports

import (
	"context"
	"encoding/json"
	"github.com/SpawNKZ/content_service/common/errors"
	"github.com/SpawNKZ/content_service/post/endpoints"
	"github.com/SpawNKZ/content_service/post/models"
	postService "github.com/SpawNKZ/content_service/post/service"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport"
	httpTransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type OperationResult struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (o *OperationResult) Succeed(data interface{}) *OperationResult {
	o.Success = true
	o.Message = ""
	o.Data = data
	return o
}

func (o *OperationResult) Error(errorMessage string) *OperationResult {
	o.Success = false
	o.Message = errorMessage
	o.Data = nil
	return o
}

var (
	opResult       = &OperationResult{}
	requiredGroups = []string{"admin"}
)

// TODO: add auth middleware
func MakeHTTPHandler(s postService.Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	options := []httpTransport.ServerOption{
		httpTransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httpTransport.ServerErrorEncoder(errorEncoder),
		httpTransport.ServerBefore(httpTransport.PopulateRequestContext, jwt.HTTPToContext()),
	}
	// GET     /programs                           retrieve programs list with all translations
	// GET     /programs/search?q=asdsadas        retrieve programs by search terms
	// GET     /programs/:id                       retrieve program by id
	// POST    /programs                         	create program
	// PUT     /programs/:id                       update the program
	// DELETE  /programs/:id                       remove the given program

	// GET     /programs/locale/:locale            retrieve programs list with specific locale
	// GET     /programs/:id/locale/:locale        retrieve program by id and locale
	// POST    /programs/locale/:locale            create program for specific locale
	// PUT     /programs/:id/locale/:locale        update the program with specific locale
	// DELETE     /programs/:id/locale/:locale        delete the program with specific locale
	sr := r.PathPrefix("/api/v1/post").Subrouter()

	sr.Methods("POST").Path("").Handler(httpTransport.NewServer(
		endpoints.MakeCreateEndpoint(s),
		decodeCreateRequest,
		encodeResponse,
		options...,
	))
	sr.Methods("GET").Path("/{id}").Handler(httpTransport.NewServer(
		endpoints.MakeGetOneEndpoint(s),
		decodeGetOneRequest,
		encodeResponse,
		options...,
	))
	sr.Methods("PUT").Path("/{id}").Handler(httpTransport.NewServer(
		endpoints.MakeUpdateEndpoint(s),
		decodeUpdateRequest,
		encodeResponse,
		options...,
	))
	sr.Methods("DELETE").Path("/{id}").Handler(httpTransport.NewServer(
		endpoints.MakeDeleteOneEndpoint(s),
		decodeDeleteOneRequest,
		encodeResponse,
		options...,
	))
	sr.Methods("GET").Path("").Handler(httpTransport.NewServer(
		endpoints.MakeGetListEndpoint(s),
		decodeGetListRequest,
		encodeResponse,
		options...,
	))
	return r
}

func decodeCreateRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req models.CreateRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeGetOneRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errors.ErrBadRouting
	}
	return models.IdRequest{ID: id}, nil
}

func decodeUpdateRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errors.ErrBadRouting
	}
	var postUpdate models.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&postUpdate); err != nil {
		return nil, err
	}

	postUpdate.ID = id

	return postUpdate, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(endpoint.Failer); ok && e.Failed() != nil {
		// Not a Go kit transports error, but a business-logic error.
		// Provide those as HTTP errors.
		errorEncoder(ctx, e.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(opResult.Succeed(response))
}

func decodeDeleteOneRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errors.ErrBadRouting
	}

	return models.IdRequest{ID: id}, nil
}

func decodeGetListRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	contentId := r.URL.Query().Get("content_id")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	return models.GetListRequest{
		ReqPagination: models.ReqPagination{
			Limit: limit, Offset: offset,
		},
		PostFilter: models.PostFilter{
			ContentId: contentId,
		},
	}, nil
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(errors.ErrorToHttpCode(err))
	json.NewEncoder(w).Encode(opResult.Error(err.Error()))
}
