package frontend

import (
	"embed"
	"gophormula/pkg/livetiming"
	"gophormula/pkg/replay"
	"gophormula/pkg/signalr"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/starfederation/datastar-go/datastar"
)

//go:embed components/* assets/*
var embedFS embed.FS

type Frontend struct {
	mux *http.ServeMux
	log *slog.Logger
	tpl *template.Template
	hub *Hub
}

func New() *Frontend {
	fe := &Frontend{
		log: slog.Default(),
		hub: NewHub(),
	}

	fe.tpl = template.Must(
		template.New("").
			Funcs(template.FuncMap{
				"urlEncode": url.QueryEscape,
			}).
			ParseFS(embedFS, "components/*.gohtml"),
	)

	assets, _ := fs.Sub(embedFS, "assets")
	fe.mux = http.NewServeMux()
	fe.mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(assets))))
	fe.mux.HandleFunc("/", fe.Index)
	fe.mux.HandleFunc("/events", fe.Events)
	fe.mux.HandleFunc("POST /replay", fe.ReplayHandler())
	fe.mux.HandleFunc("POST /live", fe.LiveHandler())

	return fe
}

func (fe *Frontend) Start(addr string) error {
	fe.log.Info("listening", "addr", addr)
	return http.ListenAndServe(addr, fe.mux)
}

func (fe *Frontend) Index(w http.ResponseWriter, r *http.Request) {
	sessions := scanSessions("data")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := fe.tpl.ExecuteTemplate(w, "index", map[string]any{
		"Sessions": sessions,
	}); err != nil {
		slog.Error("template error", "err", err)
	}
}

// Events is the persistent SSE endpoint. Datastar connects here on page load.
func (fe *Frontend) Events(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	ch := fe.hub.subscribe()
	defer fe.hub.unsubscribe(ch)

	for {
		select {
		case e, ok := <-ch:
			if !ok {
				return
			}
			if e.isScript {
				if err := sse.ExecuteScript(e.fragment); err != nil {
					return
				}
				continue
			}
			var modeOpt datastar.PatchElementOption
			if e.mode == "append" {
				modeOpt = datastar.WithModeAppend()
			} else {
				modeOpt = datastar.WithModeInner()
			}
			if err := sse.PatchElements(e.fragment, datastar.WithSelectorID(e.selector), modeOpt); err != nil {
				return
			}
		case <-r.Context().Done():
			return
		}
	}
}

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
			updater := newUpdater(fe.hub)
			var circuitFetched bool

			for msg := range ch {
				data := msg.Data()
				if data == nil {
					continue
				}
				for _, parsed := range livetiming.ParseJSON(data) {
					now := time.Now()
					updater.Apply(parsed, now)

					// Fetch circuit map once after the first SessionInfo is applied.
					if _, ok := parsed.(*livetiming.SessionInfo); ok && !circuitFetched {
						circuitFetched = true
						si := updater.Sess.Info
						year := si.StartDate.Year()
						if year == 0 {
							year = now.Year()
						}
						if cm, err := livetiming.FetchCircuitMap(si.Meeting.Circuit.Key, year); err == nil {
							bounds := boundsFromCircuitMap(cm)
							updater.SetTrack(bounds, buildTrackSVGFromMap(cm, bounds))
						} else {
							slog.Warn("live: circuit map fetch failed", "err", err)
						}
					}
				}
			}
			slog.Info("live: SignalR connection closed")
		}()
	}
}

// ReplayHandler returns an HTTP handler that starts a session replay and
// streams updates to all connected SSE clients.
func (fe *Frontend) ReplayHandler() http.HandlerFunc {
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
		if seekStr := r.URL.Query().Get("seek"); seekStr != "" {
			if d, err := time.ParseDuration(seekStr); err == nil {
				r2.SeekTo(d)
			}
		}
		bounds := r2.ScanPositionBounds()
		trackSVG := fetchAndBuildTrackSVG(path, bounds)
		ch := r2.StartAndSubscribe()

		go func() {
			fe.hub.send("active-session", "inner", filepath.Base(path))
			updater := newUpdater(fe.hub)
			updater.SetTrack(bounds, trackSVG)
			catchingUp := true

			for m := range ch {
				msg := m.(replay.Message)

				// TimingData keyframes (nil Timestamp) are the initial full-state dump;
				// skip them entirely to avoid contaminating standings.
				if _, ok := msg.Value.(*livetiming.TimingData); ok && msg.Timestamp == nil {
					continue
				}

				if msg.Catchup {
					updater.Accumulate(msg.Value)
					continue
				}

				// First real-time message: flush accumulated status and standings.
				if catchingUp {
					catchingUp = false
					updater.FlushStatus()
				}

				ts := time.Now()
				if msg.Timestamp != nil {
					ts = *msg.Timestamp
				}
				updater.Apply(msg.Value, ts)
			}
		}()

		// Return 204 so Datastar does not treat this as an SSE stream.
		w.WriteHeader(http.StatusNoContent)
	}
}

// Session represents a discoverable race session on disk.
type Session struct {
	Name string // display label (relative path from data/)
	Path string // absolute path passed to the replay engine
}

// scanSessions walks the data/ directory and returns every subdirectory
// containing an Index.json file.
func scanSessions(root string) []Session {
	var sessions []Session
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error { //nolint:errcheck
		if err != nil || !d.IsDir() {
			return nil
		}
		if _, err := os.Stat(filepath.Join(path, "Index.json")); err == nil {
			rel, _ := filepath.Rel(root, path)
			sessions = append(sessions, Session{Name: formatSessionName(rel), Path: path})
		}
		return nil
	})
	return sessions
}

// formatSessionName converts a relative session path into a human-readable label.
// e.g. "2021/2021-04-18_Emilia_Romagna_Grand_Prix/2021-04-18_Race"
//
//	→ "2021 Emilia Romagna Grand Prix Race"
func formatSessionName(rel string) string {
	parts := strings.Split(filepath.ToSlash(rel), "/")
	for i, p := range parts {
		// Strip leading YYYY-MM-DD_ date prefix if present.
		if len(p) > 11 && p[4] == '-' && p[7] == '-' && p[10] == '_' {
			p = p[11:]
		}
		parts[i] = strings.ReplaceAll(p, "_", " ")
	}
	return strings.Join(parts, " ")
}
