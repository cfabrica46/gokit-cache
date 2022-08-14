package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"cache/cmd/config"
	"cache/internal/endpoint"
	"cache/internal/entity"
	"cache/internal/service"
	"cache/internal/transport"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

func main() {
	if !config.VerifyIsDockerRun() {
		if err := config.LoadEnv(); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(os.Getenv("REDIS_HOST"))

	options := &redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: "",
		DB:       0,
	}
	db := redis.NewClient(options)

	runServer(os.Getenv("PORT"), db)
}

func runServer(port string, db *redis.Client) {
	svc := service.GetService(db)

	getGenerateTokenHandler := httptransport.NewServer(
		endpoint.MakeGenerateTokenEndpoint(svc),
		transport.DecodeRequest(entity.IDUsernameEmailSecretRequest{}),
		transport.EncodeResponse,
	)

	getExtractTokenHandler := httptransport.NewServer(
		endpoint.MakeExtractTokenEndpoint(svc),
		transport.DecodeRequest(entity.TokenSecretRequest{}),
		transport.EncodeResponse,
	)

	getSetTokenHandler := httptransport.NewServer(
		endpoint.MakeManageTokenEndpoint(svc, service.NewSetTokenState()),
		transport.DecodeRequest(entity.Token{}),
		transport.EncodeResponse,
	)

	getDeleteTokenHandler := httptransport.NewServer(
		endpoint.MakeManageTokenEndpoint(svc, service.NewDeleteTokenState()),
		transport.DecodeRequest(entity.Token{}),
		transport.EncodeResponse,
	)

	getCheckTokenHandler := httptransport.NewServer(
		endpoint.MakeCheckTokenEndpoint(svc),
		transport.DecodeRequest(entity.Token{}),
		transport.EncodeResponse,
	)

	r := mux.NewRouter()
	r.Methods(http.MethodPost).Path("/generate").Handler(getGenerateTokenHandler)
	r.Methods(http.MethodPost).Path("/extract").Handler(getExtractTokenHandler)
	r.Methods(http.MethodPost).Path("/token").Handler(getSetTokenHandler)
	r.Methods(http.MethodDelete).Path("/token").Handler(getDeleteTokenHandler)
	r.Methods(http.MethodPost).Path("/check").Handler(getCheckTokenHandler)

	log.Println("ListenAndServe on localhost:" + os.Getenv("PORT"))
	log.Println(http.ListenAndServe(":"+port, r))
}
