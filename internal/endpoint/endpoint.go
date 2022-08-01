package endpoint

import (
	"context"
	"errors"
	"fmt"

	"cache/internal/entity"
	"cache/internal/service"

	"github.com/go-kit/kit/endpoint"
)

var ErrRequest = errors.New("error to request")

// MakeGenerateTokenEndpoint ...
func MakeGenerateTokenEndpoint(svc service.Service) endpoint.Endpoint {
	return func(_ context.Context, request any) (any, error) {
		req, ok := request.(entity.IDUsernameEmailSecretRequest)
		if !ok {
			return nil, fmt.Errorf("%w: isn't of type GenerateTokenRequest", ErrRequest)
		}

		token := svc.GenerateToken(req.ID, req.Username, req.Email, []byte(req.Secret))

		return entity.Token{Token: token}, nil
	}
}

// MakeExtractTokenEndpoint ...
func MakeExtractTokenEndpoint(svc service.Service) endpoint.Endpoint {
	return func(_ context.Context, request any) (any, error) {
		var errMessage string

		req, ok := request.(entity.TokenSecretRequest)
		if !ok {
			return nil, fmt.Errorf("%w: isn't of type GenerateTokenRequest", ErrRequest)
		}

		id, username, email, err := svc.ExtractToken(req.Token, []byte(req.Secret))
		if err != nil {
			errMessage = err.Error()
		}

		return entity.IDUsernameEmailErrResponse{ID: id, Username: username, Email: email, Err: errMessage}, nil
	}
}

// MakeManageTokenEndpoint ...
func MakeManageTokenEndpoint(svc service.Service, st service.State) endpoint.Endpoint {
	return func(_ context.Context, request any) (any, error) {
		var errMessage string

		req, ok := request.(entity.Token)
		if !ok {
			return nil, fmt.Errorf("%w: isn't of type Token", ErrRequest)
		}

		err := svc.ManageToken(st, req.Token)
		if err != nil {
			errMessage = err.Error()
		}

		return entity.ErrorResponse{Err: errMessage}, nil
	}
}

// MakeCheckTokenEndpoint ...
func MakeCheckTokenEndpoint(svc service.Service) endpoint.Endpoint {
	return func(_ context.Context, request any) (any, error) {
		var errMessage string

		req, ok := request.(entity.Token)
		if !ok {
			return nil, fmt.Errorf("%w: isn't of type Token", ErrRequest)
		}

		check, err := svc.CheckToken(req.Token)
		if err != nil {
			errMessage = err.Error()
		}

		return entity.CheckErrResponse{Check: check, Err: errMessage}, nil
	}
}
