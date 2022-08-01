package transport_test

import (
	"bytes"
	"cache/internal/entity"
	"cache/internal/entity/mock"
	"cache/internal/transport"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	//nolint:gosec
	generateTokenRequestJSON = `{
		 "username": "username",
		 "email": "email@email.com",
		 "secret": "secret",
		 "id": 1
	 }`

	//nolint:gosec
	extractTokenRequestJSON = `{
		"token": "token",
		 "secret": "secret"
	}`

	//nolint:gosec
	tokenRequestJSON = `{
		"token": "token"
	}`
)

func TestDecodeRequest(t *testing.T) {
	t.Parallel()

	generateTokenReq, err := http.NewRequest(
		http.MethodPost,
		mock.URLTest,
		bytes.NewBuffer([]byte(generateTokenRequestJSON)),
	)
	if err != nil {
		assert.Error(t, err)
	}

	extractTokenReq, err := http.NewRequest(
		http.MethodPost,
		mock.URLTest,
		bytes.NewBuffer([]byte(extractTokenRequestJSON)),
	)
	if err != nil {
		assert.Error(t, err)
	}

	tokenReq, err := http.NewRequest(
		http.MethodPost,
		mock.URLTest,
		bytes.NewBuffer([]byte(tokenRequestJSON)),
	)
	if err != nil {
		assert.Error(t, err)
	}

	badReq, err := http.NewRequest(http.MethodPost, mock.URLTest, bytes.NewBuffer([]byte{}))
	if err != nil {
		assert.Error(t, err)
	}

	for _, tt := range []struct {
		inType      any
		in          *http.Request
		name        string
		outUsername string
		outEmail    string
		outToken    string
		outSecret   string
		outErr      string
		outID       int
	}{
		{
			name:        mock.NameNoError + "GenerateToken",
			inType:      entity.IDUsernameEmailSecretRequest{},
			in:          generateTokenReq,
			outID:       mock.IDTest,
			outUsername: mock.UsernameTest,
			outEmail:    mock.EmailTest,
			outSecret:   mock.SecretTest,
			outErr:      "",
		},
		{
			name:      mock.NameNoError + "ExtractToken",
			inType:    entity.TokenSecretRequest{},
			in:        extractTokenReq,
			outToken:  mock.TokenTest,
			outSecret: mock.SecretTest,
			outErr:    "",
		},
		{
			name:     mock.NameNoError + "Token",
			inType:   entity.Token{},
			in:       tokenReq,
			outToken: mock.TokenTest,
			outErr:   "",
		},
		{
			name:   "BadRequest",
			inType: entity.IDUsernameEmailSecretRequest{},
			in:     badReq,
			outErr: "EOF",
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var resultErr string

			var req any

			switch resultType := tt.inType.(type) {
			case entity.IDUsernameEmailSecretRequest:
				req, err = transport.DecodeRequest(resultType)(context.TODO(), tt.in)
				if err != nil {
					resultErr = err.Error()
				}

				result, ok := req.(entity.IDUsernameEmailSecretRequest)
				if ok {
					assert.Equal(t, tt.outID, result.ID)
					assert.Equal(t, tt.outUsername, result.Username)
					assert.Equal(t, tt.outEmail, result.Email)
					assert.Equal(t, tt.outSecret, result.Secret)
					assert.Contains(t, resultErr, tt.outErr)
				} else {
					assert.NotNil(t, err)
				}

			case entity.TokenSecretRequest:
				req, err = transport.DecodeRequest(resultType)(context.TODO(), tt.in)
				if err != nil {
					resultErr = err.Error()
				}

				result, ok := req.(entity.TokenSecretRequest)
				assert.True(t, ok)

				assert.Equal(t, tt.outToken, result.Token)
				assert.Equal(t, tt.outSecret, result.Secret)
				assert.Contains(t, resultErr, tt.outErr)

			case entity.Token:
				req, err = transport.DecodeRequest(resultType)(context.TODO(), tt.in)

				result, ok := req.(entity.Token)
				assert.True(t, ok)

				assert.Equal(t, tt.outToken, result.Token)
				assert.Contains(t, resultErr, tt.outErr)
			}
		})
	}
}

func TestEncodeResponse(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name   string
		in     any
		outErr string
	}{
		{
			name:   mock.NameNoError,
			in:     "test",
			outErr: "",
		},
		{
			name:   "ErrorEncode",
			in:     "test",
			outErr: "",
		},
		{
			name:   "ErrorBadEncode",
			in:     func() {},
			outErr: "json: unsupported type: func()",
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var resultErr string

			err := transport.EncodeResponse(context.TODO(), httptest.NewRecorder(), tt.in)
			if err != nil {
				resultErr = err.Error()
			}

			if tt.name == mock.NameNoError {
				assert.Empty(t, resultErr)
			} else {
				assert.Contains(t, resultErr, tt.outErr)
			}
		})
	}
}
