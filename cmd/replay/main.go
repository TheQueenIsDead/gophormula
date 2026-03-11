package main

import (
	"fmt"
	"gophormula/pkg/replay"
	"log/slog"
	"os"
	"path/filepath"
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
	if len(os.Args) < 2 {
		fmt.Println("Usage: replay <dir>")
	}

	dir := os.Args[1]

	replayer := replay.New()
	err := replayer.ParseGlob(filepath.Join(dir, "*"))
	if err != nil {
		slog.Error("failed to parse glob", "err", err)
		os.Exit(1)
	}
	ch := replayer.StartAndSubscribe()
	defer replayer.Close()

	for {
		select {
		case m := <-ch:
			fmt.Println(m)
		}
	}
}
