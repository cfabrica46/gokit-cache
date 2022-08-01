package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cache/internal/entity"

	httptransport "github.com/go-kit/kit/transport/http"
)

// DecodeRequest ...
func DecodeRequest[req entity.IDUsernameEmailSecretRequest |
	entity.TokenSecretRequest |
	entity.Token](request req,
) httptransport.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (any, error) {
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			return nil, fmt.Errorf("failed to decode request: %w", err)
		}

		return request, nil
	}
}

// EncodeResponse ...
func EncodeResponse(_ context.Context, w http.ResponseWriter, response any) error {
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return fmt.Errorf("failed to encode response: %w", err)
	}

	return nil
}
