package frontend

import (
	"embed"
	"gophormula/pkg/livetiming"
	"gophormula/pkg/replay"
	"gophormula/pkg/session"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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

func New(dataDir string) *Frontend {
	fe := &Frontend{
		log: slog.Default(),
		hub: NewHub(dataDir),
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
	sessions := scanSessions(fe.hub.dataDir)
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

			sess := session.New()
			catchingUp := true
			for m := range ch {
				msg := m.(replay.Message)
				// Always accumulate state even during catch-up.
				// TimingData keyframes (nil Timestamp) are skipped to avoid contaminating
				// standings with the full-state dump; all other types are always applied.
				var rerender bool
				switch msg.Value.(type) {
				case *livetiming.TimingData:
					if msg.Timestamp != nil {
						sess.Apply(msg.Value)
						rerender = true
					}
				case *livetiming.DriverList:
					sess.Apply(msg.Value)
					rerender = true
				default:
					sess.Apply(msg.Value)
				}
				// During catch-up, skip all UI updates.
				if msg.Catchup {
					continue
				}
				// First real-time message: flush accumulated status and standings.
				if catchingUp {
					catchingUp = false
					flushStatus(sess, fe.hub)
					if s := renderStandings(sess); s != "" {
						fe.hub.send("standings-panel", "inner", s)
					}
				}
				if msg.Timestamp != nil {
					fe.hub.BroadcastStatus("status-time", msg.Timestamp.Format("15:04:05"))
				}
				if pd, ok := msg.Value.(*livetiming.PositionData); ok {
					fe.hub.BroadcastCars(buildCarsSVG(pd, bounds, sess.Drivers, trackSVG))
					continue
				}
				if rerender {
					if s := renderStandings(sess); s != "" {
						fe.hub.send("standings-panel", "inner", s)
					}
				}
				if msg.Timestamp != nil {
					updateStatus(fe.hub, msg.Value)
				}
				body := formatMessage(msg.Value)
				if body == "" {
					continue
				}
				ts := "--:--:--"
				if msg.Timestamp != nil {
					ts = msg.Timestamp.Format("15:04:05")
				}
				fe.hub.Broadcast(ts, body)
			}
		}()

		// Return 204 so Datastar does not treat this as an SSE stream.
		// Sending SSE on the POST response causes Datastar to briefly drop the
		// persistent /events connection. The goroutine pushes updates through /events instead.
		w.WriteHeader(http.StatusNoContent)
	}
}

// Session represents a discoverable race session on disk.
type Session struct {
	Name string // relative path from dataDir, used as display label
	Path string // absolute path passed to the replay engine
}

// scanSessions walks dataDir and returns every directory containing Index.json.
func scanSessions(dataDir string) []Session {
	var sessions []Session
	filepath.WalkDir(dataDir, func(path string, d fs.DirEntry, err error) error { //nolint:errcheck
		if err != nil || !d.IsDir() {
			return nil
		}
		if _, err := os.Stat(filepath.Join(path, "Index.json")); err == nil {
			rel, _ := filepath.Rel(dataDir, path)
			sessions = append(sessions, Session{Name: rel, Path: path})
		}
		return nil
	})
	return sessions
}
