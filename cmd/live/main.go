package main

import (
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
	for {
		select {
		case msg := <-ch:
			for _, data := range livetiming.ParseJSON(msg.Data()) {
				switch v := data.(type) {
				case *livetiming.Heartbeat:
					slog.Debug("heartbeat", "utc", v.Utc)
				case *livetiming.CarData:
					slog.Debug("car data", "cars", len(v.Entries))
				case *livetiming.PositionData:
					slog.Debug("position data")
				case *livetiming.SessionInfo:
					slog.Info("session info", "meeting", v.Meeting.Name, "session", v.Name)
				case *livetiming.TimingData:
					slog.Debug("timing data", "drivers", len(v.Lines))
				case *livetiming.TopThree:
					slog.Debug("top three", "drivers", len(v.Lines))
				case *livetiming.TimingStats:
					slog.Debug("timing stats", "drivers", len(v.Lines))
				case *livetiming.TimingAppData:
					slog.Debug("timing app data", "drivers", len(v.Lines))
				case *livetiming.WeatherData:
					slog.Info("weather", "air", v.AirTemp, "track", v.TrackTemp)
				case *livetiming.TrackStatus:
					slog.Info("track status", "status", v.Status, "message", v.Message)
				case *livetiming.DriverList:
					slog.Info("driver list", "drivers", len(*v))
				case *livetiming.RaceControlMessages:
					slog.Info("race control", "messages", len(v.Messages))
				case *livetiming.SessionData:
					slog.Debug("session data", "points", len(v.Series))
				case *livetiming.LapCount:
					slog.Info("lap count", "current", v.CurrentLap, "total", v.TotalLaps)
				default:
					slog.Warn("unknown message type", "type", data)
				}
			}
		}
	}
}
