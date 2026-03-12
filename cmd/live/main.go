package main

import (
	"fmt"
	"gophormula/pkg/livetiming"
	"gophormula/pkg/signalr"
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
	client := signalr.NewClient(
		signalr.WithURL("https://livetiming.formula1.com/signalr"),
	)

	ch, err := client.Connect(
		[]signalr.Hub{"Streaming"},
		livetiming.AllTopics(),
	)
	if err != nil {
		slog.Error("failed to connect", "err", err)
		os.Exit(1)
	}

	slog.Info("waiting for messages")
	for msg := range ch {
		for _, data := range livetiming.ParseJSON(msg.Data()) {
			slog.Info(fmt.Sprint(data))
		}
	}
}
