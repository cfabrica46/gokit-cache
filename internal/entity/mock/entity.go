package mock

const (
	URLTest string = "localhost:8080"

	IDTest       int    = 1
	UsernameTest string = "username"
	EmailTest    string = "email@email.com"
	SecretTest   string = "secret"
	TokenTest    string = "token"

	ErrRedisClosed string = "redis: client is closed"

	NameNoError         string = "NoError"
	NameErrorRequest    string = "ErrorRequest"
	NameErrorRedisClose string = "ErrorRedisClose"
)
