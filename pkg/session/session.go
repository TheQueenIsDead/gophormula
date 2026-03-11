package session

import "gophormula/pkg/livetiming"

// Session accumulates the live state of a single F1 session from a stream of
// parsed livetiming messages. It is source-agnostic: callers feed it messages
// from a replay file or a live SignalR connection via Apply.
type Session struct {
	Info      livetiming.SessionInfo
	Drivers   map[string]livetiming.Driver
	Standings map[string]livetiming.TimingDataLine
	LapCount  livetiming.LapCount
	Weather   *livetiming.WeatherData
	Track     *livetiming.TrackStatus
	Status    *livetiming.SessionStatus
	Clock     *livetiming.ExtrapolatedClock
}

func New() *Session {
	return &Session{
		Drivers:   make(map[string]livetiming.Driver),
		Standings: make(map[string]livetiming.TimingDataLine),
	}
}

// Apply updates session state from a single parsed livetiming message.
// Returns true if any state changed.
func (s *Session) Apply(msg any) bool {
	switch v := msg.(type) {
	case *livetiming.SessionInfo:
		s.Info = *v
		return true
	case *livetiming.DriverList:
		for num, d := range *v {
			if d.Tla != "" {
				s.Drivers[num] = d
			}
		}
		return true
	case *livetiming.TimingData:
		s.mergeTimingData(v)
		return true
	case *livetiming.LapCount:
		if v.CurrentLap > 0 {
			s.LapCount.CurrentLap = v.CurrentLap
		}
		if v.TotalLaps > 0 {
			s.LapCount.TotalLaps = v.TotalLaps
		}
		return true
	case *livetiming.WeatherData:
		s.Weather = v
		return true
	case *livetiming.TrackStatus:
		s.Track = v
		return true
	case *livetiming.SessionStatus:
		s.Status = v
		return true
	case *livetiming.ExtrapolatedClock:
		s.Clock = v
		return true
	}
	return false
}

// mergeTimingData applies a delta TimingData update. Only non-zero fields
// overwrite existing values. Pointer bool fields (InPit, PitOut, Stopped,
// Retired) use nil to mean "not present in this delta", so explicit false
// values are respected and transient flags clear correctly.
func (s *Session) mergeTimingData(td *livetiming.TimingData) {
	for num, delta := range td.Lines {
		existing := s.Standings[num]
		if delta.Position != "" {
			existing.Position = delta.Position
		}
		if delta.GapToLeader != "" {
			existing.GapToLeader = delta.GapToLeader
		}
		if delta.IntervalToPositionAhead.Value != "" {
			existing.IntervalToPositionAhead = delta.IntervalToPositionAhead
		}
		if delta.NumberOfLaps > 0 {
			existing.NumberOfLaps = delta.NumberOfLaps
		}
		if delta.NumberOfPitStops > 0 {
			existing.NumberOfPitStops = delta.NumberOfPitStops
		}
		if delta.Retired != nil {
			existing.Retired = delta.Retired
		}
		if delta.PitOut != nil {
			existing.PitOut = delta.PitOut
			if pbool(delta.PitOut) {
				f := false
				existing.InPit = &f
			}
		}
		if delta.InPit != nil {
			existing.InPit = delta.InPit
			if pbool(delta.InPit) {
				f := false
				existing.PitOut = &f
			}
		}
		if delta.Stopped != nil {
			existing.Stopped = delta.Stopped
		}
		s.Standings[num] = existing
	}
}

func pbool(p *bool) bool { return p != nil && *p }
