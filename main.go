package main

import (
	"log"

	"me.zyrouge.anything_to_rss/internal/common"
	"me.zyrouge.anything_to_rss/internal/server"
)

func start() error {
	err := common.ReadEnv()
	if err != nil {
		return err
	}
	return server.StartServer()
}

func main() {
	err := start()
	if err != nil {
		log.Panicln(err)
	}
}
