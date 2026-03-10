package dash

import (
	_ "embed"
	"fmt"
	"html"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

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

// Hub manages SSE subscribers and broadcasts log lines to all connected clients.
type Hub struct {
	mu      sync.Mutex
	clients map[chan string]struct{}
	dataDir string
}

func NewHub(dataDir string) *Hub {
	return &Hub{
		clients: make(map[chan string]struct{}),
		dataDir: dataDir,
	}
}

func (h *Hub) subscribe() chan string {
	ch := make(chan string, 64)
	h.mu.Lock()
	h.clients[ch] = struct{}{}
	h.mu.Unlock()
	return ch
}

func (h *Hub) unsubscribe(ch chan string) {
	h.mu.Lock()
	delete(h.clients, ch)
	h.mu.Unlock()
	close(ch)
}

// Broadcast sends a log line to all connected SSE clients, dropping any that
// are too slow to consume (non-blocking send).
func (h *Hub) Broadcast(msg string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for ch := range h.clients {
		select {
		case ch <- msg:
		default:
		}
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

func (h *Hub) Index(w http.ResponseWriter, r *http.Request) {
	sessions := scanSessions(h.dataDir)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := pageTmpl.Execute(w, map[string]any{
		"Sessions": sessions,
		"DataDir":  h.dataDir,
	}); err != nil {
		log.Println("template error:", err)
	}
}

// Events is the persistent SSE endpoint. Datastar connects here on page load
// and receives datastar-patch-elements events prepending new log entries.
func (h *Hub) Events(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	ch := h.subscribe()
	defer h.unsubscribe(ch)

	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return
			}
			ts := time.Now().Format("15:04:05.000")
			fragment := fmt.Sprintf(
				`<div class="entry"><span class="ts">%s</span><span class="body">%s</span></div>`,
				ts, html.EscapeString(msg),
			)
			if err := sse.PatchElements(fragment, datastar.WithSelectorID("log"), datastar.WithModeAppend()); err != nil {
				return
			}
		case <-r.Context().Done():
			return
		}
	}
}

// SessionStarted sends a Datastar patch event back to the browser that initiated
// a replay, updating the active-session indicator. Call this from the /replay handler.
func SessionStarted(w http.ResponseWriter, r *http.Request, name string) {
	sse := datastar.NewSSE(w, r)
	fragment := fmt.Sprintf(`<span id="active-session">%s</span>`, html.EscapeString(name))
	if err := sse.PatchElements(fragment); err != nil {
		log.Println("SessionStarted patch error:", err)
	}
}
