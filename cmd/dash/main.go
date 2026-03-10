package main

import (
	"fmt"
	"gophormula/pkg/dash"
	"gophormula/pkg/livetiming"
	"gophormula/pkg/replay"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	dataDir := "data"
	if len(os.Args) >= 2 {
		dataDir = os.Args[1]
	}

	hub := dash.NewHub(dataDir)

	mux := http.NewServeMux()
	mux.HandleFunc("/", hub.Index)
	mux.HandleFunc("/events", hub.Events)
	mux.HandleFunc("/replay", replayHandler(hub))

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}

func replayHandler(hub *dash.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Query().Get("path")
		if path == "" {
			http.Error(w, "missing path", http.StatusBadRequest)
			return
		}

		r2 := replay.New()
		if err := r2.ParseGlob(filepath.Join(path, "*")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ch := r2.StartAndSubscribe()
		go func() {
			for m := range ch {
				msg := m.(replay.Message)
				ts := "--:--:--"
				if msg.Timestamp != nil {
					ts = msg.Timestamp.Format("15:04:05")
				}
				hub.Broadcast(ts, format(msg.Value))
			}
		}()

		dash.SessionStarted(w, r, filepath.Base(path))
	}
}

func format(msg any) string {
	switch v := msg.(type) {
	case *livetiming.Heartbeat:
		return fmt.Sprintf("Heartbeat  utc=%s", v.Utc.Format("15:04:05"))
	case *livetiming.LapCount:
		return fmt.Sprintf("LapCount   lap=%d/%d", v.CurrentLap, v.TotalLaps)
	case *livetiming.TrackStatus:
		return fmt.Sprintf("TrackStatus  status=%s  %s", v.Status, v.Message)
	case *livetiming.SessionInfo:
		return fmt.Sprintf("SessionInfo  %s — %s", v.Meeting.Name, v.Name)
	case *livetiming.WeatherData:
		return fmt.Sprintf("WeatherData  air=%s°C  track=%s°C  rain=%s", v.AirTemp, v.TrackTemp, v.Rainfall)
	case *livetiming.RaceControlMessages:
		if len(v.Messages) > 0 {
			last := v.Messages[fmt.Sprintf("%d", len(v.Messages)-1)]
			return fmt.Sprintf("RaceControl  [%s] %s", last.Category, last.Message)
		}
		return "RaceControl"
	case *livetiming.TimingData:
		return fmt.Sprintf("TimingData   %d drivers", len(v.Lines))
	case *livetiming.DriverList:
		return fmt.Sprintf("DriverList   %d drivers", len(*v))
	case *livetiming.ExtrapolatedClock:
		return fmt.Sprintf("Clock  remaining=%s  extrapolating=%v", v.Remaining, v.Extrapolating)
	default:
		return fmt.Sprintf("%T", msg)
	}
}
