package main

import (
	"context"
	"log"
	"os"
	"strings"

	"benkyou/server"
	"benkyou/service"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kanjiDir := os.Getenv("KANJI_DIR")
	if kanjiDir == "" {
		log.Fatal("empty kanji dir")
	}

	kanjiSvc, err := service.NewKanjiService(kanjiDir)
	if err != nil {
		log.Fatal(err)
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		log.Fatal("empty HTTP port")
	}

	if !strings.Contains(httpPort, ":") {
		log.Fatal("invalid http port format")
	}

	server.RunServer(ctx, kanjiSvc, httpPort)
}
