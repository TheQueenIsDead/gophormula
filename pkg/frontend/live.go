package frontend

import (
	"gophormula/pkg/livetiming"
	"gophormula/pkg/replay"
	"gophormula/pkg/session"
	"gophormula/pkg/signalr"
	"log/slog"
	"net/http"
	"time"
)

// LiveHandler connects to the F1 live timing SignalR feed and streams updates
// to all connected SSE clients.
func (fe *Frontend) LiveHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client := signalr.NewClient(
			signalr.WithURL("https://livetiming.formula1.com/signalr"),
		)
		ch, err := client.Connect([]signalr.Hub{"Streaming"}, livetiming.AllTopics())
		if err != nil {
			slog.Error("live: SignalR connect failed", "err", err)
			http.Error(w, "could not connect to F1 live timing: "+err.Error(), http.StatusBadGateway)
			return
		}
		slog.Info("live: connected to F1 SignalR")

		// Respond with 204 so Datastar does not treat this as an SSE stream.
		// Sending SSE on the POST response causes Datastar to briefly drop the
		// persistent /events connection, resulting in missed snapshot updates.
		// The goroutine pushes the active-session label through /events instead.
		w.WriteHeader(http.StatusNoContent)

		go func() {
			fe.hub.send("active-session", "inner", "Live")

			sess := session.New()
			var bounds replay.PositionBounds
			var trackSVG string

			for msg := range ch {
				data := msg.Data()
				if data == nil {
					continue
				}
				results := livetiming.ParseJSON(data)
				for _, parsed := range results {
					now := time.Now()

					// Position data: render car dots immediately.
					if pd, ok := parsed.(*livetiming.PositionData); ok {
						fe.hub.BroadcastCars(buildCarsSVG(pd, bounds, sess.Drivers, trackSVG))
						fe.hub.BroadcastStatus("status-time", now.Format("15:04:05"))
						continue
					}

					// SessionInfo: fetch circuit map once to get track SVG and bounds.
					if _, ok := parsed.(*livetiming.SessionInfo); ok {
						sess.Apply(parsed)
						if trackSVG == "" {
							si := sess.Info
							year := si.StartDate.Year()
							if year == 0 {
								year = now.Year()
							}
							if cm, err := livetiming.FetchCircuitMap(si.Meeting.Circuit.Key, year); err == nil {
								bounds = boundsFromCircuitMap(cm)
								trackSVG = buildTrackSVGFromMap(cm, bounds)
							} else {
								slog.Warn("live: circuit map fetch failed", "err", err)
							}
						}
						updateStatus(fe.hub, parsed)
						if body := formatMessage(parsed); body != "" {
							fe.hub.Broadcast(now.Format("15:04:05"), body)
						}
						continue
					}

					var rerender bool
					switch parsed.(type) {
					case *livetiming.TimingData, *livetiming.DriverList:
						rerender = sess.Apply(parsed)
					default:
						sess.Apply(parsed)
					}

					if rerender {
						if s := renderStandings(sess); s != "" {
							fe.hub.send("standings-panel", "inner", s)
						}
					}
					updateStatus(fe.hub, parsed)
					fe.hub.BroadcastStatus("status-time", now.Format("15:04:05"))

					if body := formatMessage(parsed); body != "" {
						fe.hub.Broadcast(now.Format("15:04:05"), body)
					}
				}
			}
			slog.Info("live: SignalR connection closed")
		}()
	}
}
