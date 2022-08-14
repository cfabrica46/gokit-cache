package config

import (
	"fmt"
	"os"
)

const (
	DBDriver = "postgres"
)

func VerifyIsDockerRun() (check bool) {
	isDocker := os.Getenv("DOCKER")

	return isDocker == "true"
}

func LoadEnv() (err error) {
	err = os.Setenv("PORT", "9090")
	if err != nil {
		return fmt.Errorf("error to set env:%w", err)
	}

	err = os.Setenv("REDIS_HOST", "localhost")
	if err != nil {
		return fmt.Errorf("error to set env:%w", err)
	}

	err = os.Setenv("REDIS_PORT", "6379")
	if err != nil {
		return fmt.Errorf("error to set env:%w", err)
	}

	return nil
}
