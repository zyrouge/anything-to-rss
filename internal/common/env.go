package common

import (
	"errors"
	"fmt"
	"os"
)

const HTTP_LISTEN_ADDR_ENV = "HTTP_LISTEN_ADDR"
const GITHUB_CONTAINER_REGISTRY_API_TOKEN_ENV = "GITHUB_CONTAINER_REGISTRY_API_TOKEN"

type Env struct {
	HttpListenAddr                  string
	GitHubContainerRegistryApiToken string
}

var env *Env

func GetEnv() (*Env, error) {
	if env == nil {
		return nil, errors.New("Env is not set")
	}
	return env, nil
}

func ReadEnv() error {
	current := Env{
		HttpListenAddr:                  os.Getenv(HTTP_LISTEN_ADDR_ENV),
		GitHubContainerRegistryApiToken: os.Getenv(GITHUB_CONTAINER_REGISTRY_API_TOKEN_ENV),
	}
	if current.HttpListenAddr == "" {
		return fmt.Errorf("missing env: %s", HTTP_LISTEN_ADDR_ENV)
	}
	if current.GitHubContainerRegistryApiToken == "" {
		return fmt.Errorf("missing env: %s", GITHUB_CONTAINER_REGISTRY_API_TOKEN_ENV)
	}
	env = &current
	return nil
}
