package frontend

import (
	"fmt"
	"html"
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

// BroadcastCars replaces the inner HTML of #plot-panel with a full SVG.
// Targeting the HTML parent (not the SVG element itself) ensures the browser
// parses car elements in the correct SVG namespace context. Idiomorph diffs
// the SVG children by id so CSS transitions can animate cx/cy in place.
func (h *Hub) BroadcastCars(fragment string) {
	h.send("plot-panel", "inner", fragment)
}

// BroadcastTrack sends the circuit outline SVG fragment to #track-outline,
// used for the initial track render before any position data arrives.
func (h *Hub) BroadcastTrack(fragment string) {
	h.send("track-outline", "inner", fragment)
}

// BroadcastStatus replaces the inner HTML of the named status element.
func (h *Hub) BroadcastStatus(id, fragment string) {
	h.send(id, "inner", fragment)
}
