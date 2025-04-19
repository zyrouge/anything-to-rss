package server

import (
	"log"
	"net"
	"net/http"

	"me.zyrouge.anything_to_rss/internal/common"
	"me.zyrouge.anything_to_rss/internal/sources"
)

func StartServer() error {
	env, err := common.GetEnv()
	if err != nil {
		return err
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/docker-hub-tags", sources.RouteDockerHubTags)
	mux.HandleFunc("/github-container-registry-versions", sources.RouteGitHubContainerRegistryVersions)
	server := http.Server{
		Handler: mux,
	}
	log.Printf("Listening on %s\n", env.HttpListenAddr)
	listener, err := net.Listen("tcp", env.HttpListenAddr)
	if err != nil {
		return err
	}
	return server.Serve(listener)
}
