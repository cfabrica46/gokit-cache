package endpoint_test

import (
	"context"
	"testing"

	"cache/internal/endpoint"
	"cache/internal/entity"
	"cache/internal/entity/mock"
	"cache/internal/service"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type incorrectRequest struct {
	incorrect bool
}

func TestMakeGenerateTokenEndpoint(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		in     any
		name   string
		outErr string
	}{
		{
			name: mock.NameNoError,
			in: entity.IDUsernameEmailSecretRequest{
				ID:       mock.IDTest,
				Username: mock.UsernameTest,
				Email:    mock.EmailTest,
				Secret:   mock.SecretTest,
			},
			outErr: "",
		},
		{
			name: mock.NameErrorRequest,
			in: incorrectRequest{
				incorrect: true,
			},
			outErr: "isn't of type",
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var resultErr string

			mr, err := miniredis.Run()
			if err != nil {
				assert.Error(t, err)
			}

			client := redis.NewClient(&redis.Options{Addr: mr.Addr()})

			svc := service.GetService(client)

			r, err := endpoint.MakeGenerateTokenEndpoint(svc)(context.TODO(), tt.in)
			if err != nil {
				resultErr = err.Error()
			}

			result, ok := r.(entity.Token)
			if !ok {
				if tt.name != mock.NameErrorRequest {
					assert.Fail(t, "response is not of the type indicated")
				}
			}

			if tt.name == mock.NameNoError {
				assert.Empty(t, resultErr)
				assert.NotEmpty(t, result.Token)
			} else {
				assert.Contains(t, resultErr, tt.outErr)
				assert.Empty(t, result.Token)
			}
		})
	}
}

func TestMakeExtractTokenEndpoint(t *testing.T) {
	t.Parallel()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       mock.IDTest,
		"username": mock.UsernameTest,
		"email":    mock.EmailTest,
		"uuid":     uuid.NewString(),
	})

	tokenSigned, err := token.SignedString([]byte(mock.SecretTest))
	if err != nil {
		assert.Error(t, err)
	}

	for _, tt := range []struct {
		name   string
		in     any
		outErr string
	}{
		{
			name: mock.NameNoError,
			in: entity.TokenSecretRequest{
				Token:  tokenSigned,
				Secret: mock.SecretTest,
			},
			outErr: "",
		},
		{
			name: mock.NameErrorRequest,
			in: incorrectRequest{
				incorrect: true,
			},
			outErr: "isn't of type",
		},
		{
			name: "ErrorNotValidToken",
			in: entity.TokenSecretRequest{
				Token:  "",
				Secret: mock.SecretTest,
			},
			outErr: "token contains an invalid number of segments",
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var resultErr string
			var mr *miniredis.Miniredis
			var req any

			mr, err = miniredis.Run()
			if err != nil {
				assert.Error(t, err)
			}

			client := redis.NewClient(&redis.Options{Addr: mr.Addr()})

			svc := service.GetService(client)

			req, err = endpoint.MakeExtractTokenEndpoint(svc)(context.TODO(), tt.in)
			if err != nil {
				resultErr = err.Error()
			}

			result, ok := req.(entity.IDUsernameEmailErrResponse)
			if !ok {
				if tt.name != mock.NameErrorRequest {
					assert.Fail(t, "response is not of the type indicated")
				}
			} else {
				resultErr = result.Err
			}

			if tt.name == mock.NameNoError {
				assert.Empty(t, result.Err)
			} else {
				assert.Contains(t, resultErr, tt.outErr)
			}
		})
	}
}

func TestMakeManageTokenEndpoint(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		inState service.State
		in      any
		name    string
		outErr  string
	}{
		{
			name:    mock.NameNoError,
			in:      entity.Token{Token: "token"},
			inState: service.NewSetTokenState(),
			outErr:  "",
		},
		{
			name: mock.NameErrorRequest,
			in: incorrectRequest{
				incorrect: true,
			},
			outErr: "isn't of type",
		},
		{
			name:    mock.NameErrorRedisClose,
			in:      entity.Token{Token: ""},
			inState: service.NewSetTokenState(),
			outErr:  mock.ErrRedisClosed,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var resultErr string

			mr, err := miniredis.Run()
			if err != nil {
				assert.Error(t, err)
			}

			client := redis.NewClient(&redis.Options{Addr: mr.Addr()})

			svc := service.GetService(client)

			if tt.name == mock.NameErrorRedisClose {
				svc.DB.Close()
			}

			r, err := endpoint.MakeManageTokenEndpoint(svc, tt.inState)(context.TODO(), tt.in)
			if err != nil {
				resultErr = err.Error()
			}

			result, ok := r.(entity.ErrorResponse)
			if !ok {
				if tt.name != mock.NameErrorRequest {
					assert.Fail(t, "response is not of the type indicated")
				}
			} else {
				resultErr = result.Err
			}

			if tt.name == mock.NameNoError {
				assert.Empty(t, result.Err)
			} else {
				assert.Contains(t, resultErr, tt.outErr)
			}
		})
	}
}

func TestMakeCheckTokenEndpoint(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name   string
		in     any
		outErr string
	}{
		{
			name:   mock.NameNoError,
			in:     entity.Token{Token: "token"},
			outErr: "",
		},
		{
			name: mock.NameErrorRequest,
			in: incorrectRequest{
				incorrect: true,
			},
			outErr: "isn't of type",
		},
		{
			name:   mock.NameErrorRedisClose,
			in:     entity.Token{Token: ""},
			outErr: mock.ErrRedisClosed,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var resultErr string

			mr, err := miniredis.Run()
			if err != nil {
				assert.Error(t, err)
			}

			client := redis.NewClient(&redis.Options{Addr: mr.Addr()})

			svc := service.GetService(client)

			if tt.name == mock.NameErrorRedisClose {
				svc.DB.Close()
			}

			r, err := endpoint.MakeCheckTokenEndpoint(svc)(context.TODO(), tt.in)
			if err != nil {
				resultErr = err.Error()
			}

			result, ok := r.(entity.CheckErrResponse)
			if !ok {
				if tt.name != mock.NameErrorRequest {
					assert.Fail(t, "response is not of the type indicated")
				}
			} else {
				resultErr = result.Err
			}

			if tt.name == mock.NameNoError {
				assert.Empty(t, result.Err)
			} else {
				assert.Contains(t, resultErr, tt.outErr)
			}
		})
	}
}
