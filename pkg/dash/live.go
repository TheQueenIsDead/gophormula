package dash

import (
	"gophormula/pkg/livetiming"
	"gophormula/pkg/replay"
	"gophormula/pkg/session"
	"gophormula/pkg/signalr"
	"log/slog"
	"math"
	"net/http"
	"time"
)

// LiveHandler connects to the F1 live timing SignalR feed and streams updates
// to all connected SSE clients.
func (h *Hub) LiveHandler() http.HandlerFunc {
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

		go func() {
			sess := session.New()
			var bounds replay.PositionBounds
			var trackSVG string

			for msg := range ch {
				data := msg.Data()
				slog.Debug("live: raw message", "raw", string(msg.Raw), "data_len", len(data))
				if data == nil {
					continue
				}
				results := livetiming.ParseJSON(data)
				slog.Debug("live: parsed", "count", len(results))
				for _, parsed := range results {
					now := time.Now()

					// Position data: render car dots immediately.
					if pd, ok := parsed.(*livetiming.PositionData); ok {
						h.BroadcastCars(buildCarsSVG(pd, bounds, sess.Drivers, trackSVG))
						h.BroadcastStatus("status-time", now.Format("15:04:05"))
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
						updateStatus(h, parsed)
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
						s := renderStandings(sess)
						slog.Info("live: standings update", "standings_entries", len(sess.Standings), "html_len", len(s), "clients", len(h.clients))
						if s != "" {
							h.send("standings-panel", "inner", s)
						}
					}
					updateStatus(h, parsed)
					h.BroadcastStatus("status-time", now.Format("15:04:05"))

					if body := formatMessage(parsed); body != "" {
						h.Broadcast(now.Format("15:04:05"), body)
					}
				}
			}
			slog.Info("live: SignalR connection closed")
		}()

		SessionStarted(w, r, "Live")
	}
}

// boundsFromCircuitMap derives PositionBounds from a Multiviewer circuit map.
// The circuit map and F1 position data share the same coordinate space, so the
// circuit extent can normalise live car positions onto the SVG canvas.
func boundsFromCircuitMap(cm *livetiming.CircuitMap) replay.PositionBounds {
	if cm == nil || len(cm.X) == 0 {
		return replay.PositionBounds{}
	}
	minX, maxX := math.MaxInt, math.MinInt
	minY, maxY := math.MaxInt, math.MinInt
	for i := range cm.X {
		x, y := int(math.Round(cm.X[i])), int(math.Round(cm.Y[i]))
		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
	}
	return replay.PositionBounds{MinX: minX, MaxX: maxX, MinY: minY, MaxY: maxY, Valid: true}
}
