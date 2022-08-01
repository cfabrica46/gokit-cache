package service

import (
	"errors"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

const (
	lifeOfToken int = 10
)

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrClaims                  = errors.New("error to claims")
)

type Service interface {
	GenerateToken(int, string, string, []byte) string
	ExtractToken(string, []byte) (int, string, string, error)
	ManageToken(State, string) error
	CheckToken(string) (bool, error)
}

// service ...
type service struct {
	DB *redis.Client
}

// GetService ...
func GetService(db *redis.Client) *service {
	return &service{db}
}

// GenerateToken ...
func (service) GenerateToken(id int, username, email string, secret []byte) (token string) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       id,
		"username": username,
		"email":    email,
		"uuid":     uuid.NewString(),
	})

	token, _ = t.SignedString(secret)

	return token
}

// ExtractToken ...
func (service) ExtractToken(token string, secret []byte) (id int, username, email string, err error) {
	t, err := jwt.Parse(token, KeyFunc(secret))
	if err != nil {
		return 0, "", "", fmt.Errorf("error to extract token: %w", err)
	}

	claims, _ := t.Claims.(jwt.MapClaims)

	idAux, ok := claims["id"].(float64)
	if !ok {
		return 0, "", "", fmt.Errorf("%w: claims['id'] isn't of type float64", ErrClaims)
	}

	id = int(idAux)

	username, ok = claims["username"].(string)
	if !ok {
		return 0, "", "", fmt.Errorf("%w: claims['username'] isn't of type string", ErrClaims)
	}

	email, ok = claims["email"].(string)
	if !ok {
		return 0, "", "", fmt.Errorf("%w: claims['email'] isn't of type string", ErrClaims)
	}

	return id, username, email, nil
}

// ManageToken ...
func (s *service) ManageToken(st State, token string) (err error) {
	err = st.ManageToken(s.DB, token)
	if err != nil {
		return fmt.Errorf("error when managing token: %w", err)
	}

	return nil
}

// CheckToken ...
func (s service) CheckToken(token string) (check bool, err error) {
	result, err := s.DB.Get(token).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}

		return false, fmt.Errorf("error to get token: %w", err)
	}

	if result == "1" {
		check = true
	}

	return check, nil
}

func KeyFunc(secret []byte) func(token *jwt.Token) (any, error) {
	return func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}

		return secret, nil
	}
}
