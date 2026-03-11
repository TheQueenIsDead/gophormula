package main

import (
	"gophormula/pkg/dash"
	"log/slog"
	"net/http"
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
	dataDir := "data"
	if len(os.Args) >= 2 {
		dataDir = os.Args[1]
	}

	hub := dash.NewHub(dataDir)

	mux := http.NewServeMux()
	mux.HandleFunc("/", hub.Index)
	mux.HandleFunc("/events", hub.Events)
	mux.HandleFunc("/replay", hub.ReplayHandler())
	mux.HandleFunc("/live", hub.LiveHandler())

	slog.Info("listening", "addr", ":1234")
	if err := http.ListenAndServe(":1234", mux); err != nil {
		slog.Error("server error", "err", err)
		os.Exit(1)
	}
}
