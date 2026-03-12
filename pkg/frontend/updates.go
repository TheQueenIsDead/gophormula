package frontend

import (
	"gophormula/pkg/livetiming"
	"gophormula/pkg/replay"
	"gophormula/pkg/session"
	"time"
)

// Updater holds session state and pushes DOM patches to connected clients
// as messages arrive from either the replay engine or the live SignalR feed.
type Updater struct {
	hub      *Hub
	Sess     *session.Session
	bounds   replay.PositionBounds
	trackSVG string
}

func newUpdater(hub *Hub) *Updater {
	return &Updater{
		hub:  hub,
		Sess: session.New(),
	}
}

// SetTrack updates the circuit bounds and broadcasts the track outline to all clients.
func (u *Updater) SetTrack(bounds replay.PositionBounds, svg string) {
	u.bounds = bounds
	u.trackSVG = svg
	u.hub.BroadcastTrack(svg)
}

// Accumulate applies value to session state without pushing any UI updates.
// Used during replay seek catchup to build up session state before rendering begins.
func (u *Updater) Accumulate(value any) {
	u.Sess.Apply(value)
}

// FlushStatus pushes all accumulated session state to connected clients.
// Called once after replay seek/fast-forward catchup completes.
func (u *Updater) FlushStatus() {
	flushStatus(u.Sess, u.hub)
	if s := renderStandings(u.Sess); s != "" {
		u.hub.send("standings-panel", "inner", s)
	}
}

// Apply processes a single parsed livetiming message at ts, updating session
// state and pushing the appropriate DOM patches to all connected clients.
func (u *Updater) Apply(value any, ts time.Time) {
	// Position data: render car dots immediately, no log entry.
	if pd, ok := value.(*livetiming.PositionData); ok {
		for _, frame := range pd.Position {
			u.hub.BroadcastCars(buildCarsSVG(frame, u.bounds, u.Sess.Drivers, u.trackSVG))
		}
		u.hub.BroadcastStatus("status-time", ts.Format("15:04:05"))
		return
	}

	if changed := u.Sess.Apply(value); changed {
		if s := renderStandings(u.Sess); s != "" {
			u.hub.send("standings-panel", "inner", s)
		}
	}
	updateStatus(u.hub, value)
	u.hub.BroadcastStatus("status-time", ts.Format("15:04:05"))

	if body := formatMessage(value); body != "" {
		u.hub.Broadcast(ts.Format("15:04:05"), body)
	}
}
