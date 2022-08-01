package service

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type State interface {
	ManageToken(*redis.Client, string) error
}

type (
	SetTokenState    struct{}
	DeleteTokenState struct{}
)

func NewSetTokenState() SetTokenState {
	return SetTokenState{}
}

func (SetTokenState) ManageToken(db *redis.Client, token string) (err error) {
	err = db.Set(token, true, time.Minute*time.Duration(lifeOfToken)).Err()
	if err != nil {
		return fmt.Errorf("error to set token: %w", err)
	}

	return nil
}

func NewDeleteTokenState() DeleteTokenState {
	return DeleteTokenState{}
}

func (DeleteTokenState) ManageToken(db *redis.Client, token string) (err error) {
	if err = db.Del(token).Err(); err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	return nil
}
