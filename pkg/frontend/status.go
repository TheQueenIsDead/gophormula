package frontend

import (
	"fmt"
	"gophormula/pkg/session"
	"html"
)

// flushStatus broadcasts the accumulated status values from sess to all
// connected clients. Called once after seek/fast-forward catchup completes.
func flushStatus(sess *session.Session, h *Hub) {
	if sess.LapCount.CurrentLap > 0 || sess.LapCount.TotalLaps > 0 {
		h.BroadcastStatus("status-lap", fmt.Sprintf("Lap %d / %d", sess.LapCount.CurrentLap, sess.LapCount.TotalLaps))
	}
	if sess.Weather != nil {
		h.BroadcastStatus("status-weather",
			fmt.Sprintf("Air %s°C · Track %s°C · Wind %skm/h · Rain %s",
				sess.Weather.AirTemp, sess.Weather.TrackTemp, sess.Weather.WindSpeed, sess.Weather.Rainfall))
	}
	if sess.Track != nil {
		color := trackStatusColor(sess.Track.Status)
		h.BroadcastStatus("status-track",
			fmt.Sprintf(`<span style="color:%s;font-weight:bold">%s</span>`, color, html.EscapeString(sess.Track.Message)))
	}
	if sess.Status != nil {
		color := sessionStatusColor(sess.Status.Status)
		h.BroadcastStatus("status-session",
			fmt.Sprintf(`<span class="status-dot" style="background:%s" data-tip="%s"></span>`,
				color, html.EscapeString(sess.Status.Status)))
	}
	if sess.Clock != nil {
		h.BroadcastStatus("status-remaining", html.EscapeString(sess.Clock.Remaining))
	}
}
