package dash

import (
	_ "embed"
	"fmt"
	"html"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"

	"github.com/starfederation/datastar-go/datastar"
)

//go:embed page.html
var pageTmplSrc string

var pageTmpl = template.Must(
	template.New("page").
		Funcs(template.FuncMap{
			"urlEncode": url.QueryEscape,
		}).
		Parse(pageTmplSrc),
)

// Hub manages SSE subscribers and broadcasts updates to all connected clients.
type Hub struct {
	mu      sync.Mutex
	clients map[chan entry]struct{}
	dataDir string
}

func NewHub(dataDir string) *Hub {
	return &Hub{
		clients: make(map[chan entry]struct{}),
		dataDir: dataDir,
	}
}

func (h *Hub) subscribe() chan entry {
	ch := make(chan entry, 64)
	h.mu.Lock()
	h.clients[ch] = struct{}{}
	h.mu.Unlock()
	return ch
}

func (h *Hub) unsubscribe(ch chan entry) {
	h.mu.Lock()
	delete(h.clients, ch)
	h.mu.Unlock()
	close(ch)
}

// entry is a generic DOM patch: replace the inner HTML of selector with fragment,
// or append to it when mode is "append".
type entry struct {
	selector string
	mode     string // "append" or "inner"
	fragment string
}

func (h *Hub) send(selector, mode, fragment string) {
	if fragment == "" {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	for ch := range h.clients {
		select {
		case ch <- entry{selector: selector, mode: mode, fragment: fragment}:
		default:
		}
	}
}

// Broadcast appends a timestamped log entry to #log.
func (h *Hub) Broadcast(ts, body string) {
	fragment := fmt.Sprintf(
		`<div class="entry"><span class="ts">%s</span><span class="body">%s</span></div>`,
		html.EscapeString(ts), html.EscapeString(body),
	)
	h.send("log", "append", fragment)
}

// BroadcastCars replaces #plot-panel with a fresh SVG (avoids SVG namespace issues with innerHTML).
func (h *Hub) BroadcastCars(fragment string) {
	h.send("plot-panel", "inner", fragment)
}

// BroadcastStatus replaces the inner HTML of the named status element.
func (h *Hub) BroadcastStatus(id, fragment string) {
	h.send(id, "inner", fragment)
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

func (h *Hub) Index(w http.ResponseWriter, r *http.Request) {
	sessions := scanSessions(h.dataDir)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := pageTmpl.Execute(w, map[string]any{
		"Sessions": sessions,
		"DataDir":  h.dataDir,
	}); err != nil {
		slog.Error("template error", "err", err)
	}
}

// Events is the persistent SSE endpoint. Datastar connects here on page load.
func (h *Hub) Events(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	ch := h.subscribe()
	defer h.unsubscribe(ch)

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

// SessionStarted sends a Datastar patch event back to the browser that initiated
// a replay, updating the active-session indicator.
func SessionStarted(w http.ResponseWriter, r *http.Request, name string) {
	sse := datastar.NewSSE(w, r)
	fragment := fmt.Sprintf(`<span id="active-session">%s</span>`, html.EscapeString(name))
	if err := sse.PatchElements(fragment); err != nil {
		slog.Warn("SessionStarted patch error", "err", err)
	}
}
