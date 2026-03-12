package main

import (
	"gophormula/pkg/frontend"
	"log"
	"log/slog"
	"os"
)

func initLogging() {
	level := slog.LevelInfo
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		_ = level.UnmarshalText([]byte(v))
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})))
}

func main() {
	initLogging()
	if err := frontend.New().Start(":1234"); err != nil {
		log.Fatal(err)
	}
}
