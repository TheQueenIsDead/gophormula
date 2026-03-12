package frontend

import (
	"fmt"
	"html"
	"strings"
	"sync"
)

// Hub manages SSE subscribers and broadcasts updates to all connected clients.
type Hub struct {
	mu      sync.Mutex
	clients map[chan entry]struct{}
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[chan entry]struct{}),
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

// entry is either a DOM patch (selector+mode+fragment) or a script execution (isScript=true, fragment=script).
type entry struct {
	selector string
	mode     string // "append" or "inner"
	fragment string
	isScript bool
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

// BroadcastScript executes a JavaScript snippet on all connected clients via
// Datastar's datastar-execute-script SSE event.
func (h *Hub) BroadcastScript(script string) {
	if script == "" {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	for ch := range h.clients {
		select {
		case ch <- entry{fragment: script, isScript: true}:
		default:
		}
	}
}

// BroadcastTrack sets the circuit outline on #track-outline via script injection.
// Direct innerHTML assignment on an SVG element preserves SVG namespace context,
// whereas Datastar's PatchElements parses fragments as HTML (losing the namespace).
func (h *Hub) BroadcastTrack(fragment string) {
	if fragment == "" {
		return
	}
	escaped := strings.ReplaceAll(fragment, "\\", "\\\\")
	escaped = strings.ReplaceAll(escaped, "`", "\\`")
	h.BroadcastScript("document.getElementById('track-outline').innerHTML=`" + escaped + "`")
}

// BroadcastStatus replaces the inner HTML of the named status element.
func (h *Hub) BroadcastStatus(id, fragment string) {
	h.send(id, "inner", fragment)
}
