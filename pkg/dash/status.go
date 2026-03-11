package dash

import (
	"fmt"
	"gophormula/pkg/livetiming"
	"html"
)

// statusState accumulates the latest value of each status-bar element so that
// after a seek/fast-forward the correct state can be flushed to the UI on the
// first real-time message.
type statusState struct {
	lapCount      livetiming.LapCount
	weather       *livetiming.WeatherData
	trackStatus   *livetiming.TrackStatus
	sessionStatus *livetiming.SessionStatus
	clock         *livetiming.ExtrapolatedClock
}

func newStatusState() *statusState { return &statusState{} }

// merge updates the accumulator from any message type that affects the status bar.
func (s *statusState) merge(msg any) {
	switch v := msg.(type) {
	case *livetiming.LapCount:
		if v.CurrentLap > 0 {
			s.lapCount.CurrentLap = v.CurrentLap
		}
		if v.TotalLaps > 0 {
			s.lapCount.TotalLaps = v.TotalLaps
		}
	case *livetiming.WeatherData:
		s.weather = v
	case *livetiming.TrackStatus:
		s.trackStatus = v
	case *livetiming.SessionStatus:
		s.sessionStatus = v
	case *livetiming.ExtrapolatedClock:
		s.clock = v
	}
}

// flush broadcasts the accumulated status values to all connected clients.
func (s *statusState) flush(h *Hub) {
	if s.lapCount.CurrentLap > 0 || s.lapCount.TotalLaps > 0 {
		h.BroadcastStatus("status-lap", fmt.Sprintf("Lap %d / %d", s.lapCount.CurrentLap, s.lapCount.TotalLaps))
	}
	if s.weather != nil {
		h.BroadcastStatus("status-weather",
			fmt.Sprintf("Air %s°C · Track %s°C · Wind %skm/h · Rain %s",
				s.weather.AirTemp, s.weather.TrackTemp, s.weather.WindSpeed, s.weather.Rainfall))
	}
	if s.trackStatus != nil {
		color := trackStatusColor(s.trackStatus.Status)
		h.BroadcastStatus("status-track",
			fmt.Sprintf(`<span style="color:%s;font-weight:bold">%s</span>`, color, html.EscapeString(s.trackStatus.Message)))
	}
	if s.sessionStatus != nil {
		color := sessionStatusColor(s.sessionStatus.Status)
		h.BroadcastStatus("status-session",
			fmt.Sprintf(`<span class="status-dot" style="background:%s" data-tip="%s"></span>`,
				color, html.EscapeString(s.sessionStatus.Status)))
	}
	if s.clock != nil {
		h.BroadcastStatus("status-remaining", html.EscapeString(s.clock.Remaining))
	}
}