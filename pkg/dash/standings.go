package dash

import (
	"fmt"
	"gophormula/pkg/livetiming"
	"html"
	"sort"
	"strconv"
	"strings"
)

// standingsState accumulates incremental TimingData and DriverList deltas so
// we can always render the full current standings.
type standingsState struct {
	lines   map[string]livetiming.TimingDataLine
	drivers map[string]livetiming.Driver
}

func newStandingsState() *standingsState {
	return &standingsState{
		lines:   make(map[string]livetiming.TimingDataLine),
		drivers: make(map[string]livetiming.Driver),
	}
}

func pbool(p *bool) bool { return p != nil && *p }

// mergeTimingData applies a delta TimingData update to the accumulated state.
// Pointer bool fields (InPit, PitOut, Stopped, Retired) use nil to mean
// "not present in this delta", so both true and false values are respected.
func (s *standingsState) mergeTimingData(td *livetiming.TimingData) {
	for num, delta := range td.Lines {
		existing := s.lines[num]
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
		s.lines[num] = existing
	}
}

// mergeDriverList updates the driver info map from a DriverList message.
func (s *standingsState) mergeDriverList(dl *livetiming.DriverList) {
	for num, d := range *dl {
		if d.Tla != "" {
			s.drivers[num] = d
		}
	}
}

// render builds the inner HTML for the standings panel, sorted by race position.
func (s *standingsState) render() string {
	if len(s.lines) == 0 {
		return ""
	}

	type row struct {
		num  string
		line livetiming.TimingDataLine
		pos  int
	}
	rows := make([]row, 0, len(s.lines))
	for num, line := range s.lines {
		pos, _ := strconv.Atoi(line.Position)
		rows = append(rows, row{num: num, line: line, pos: pos})
	}
	// Retired drivers go last; within retired group, keep position order.
	sort.Slice(rows, func(i, j int) bool {
		ri, rj := pbool(rows[i].line.Retired), pbool(rows[j].line.Retired)
		if ri != rj {
			return !ri
		}
		if rows[i].pos != rows[j].pos {
			return rows[i].pos < rows[j].pos
		}
		return rows[i].num < rows[j].num
	})

	var sb strings.Builder
	for _, r := range rows {
		driver := s.drivers[r.num]
		tla := driver.Tla
		if tla == "" {
			tla = "#" + r.num
		}
		colour := "#888888"
		if driver.TeamColour != "" {
			colour = "#" + driver.TeamColour
		}

		gap := html.EscapeString(r.line.GapToLeader)
		if r.pos == 1 {
			gap = "Leader"
		}

		retired := pbool(r.line.Retired)
		rowStyle := fmt.Sprintf("border-left:3px solid %s", colour)
		rowClass := "sr"
		if retired {
			rowClass += " sr-out"
		}

		badge := ""
		switch {
		case retired:
			badge = `<span class="sr-badge sr-badge-out">OUT</span>`
		case pbool(r.line.InPit):
			badge = `<span class="sr-badge sr-badge-pit">PIT</span>`
		case pbool(r.line.PitOut):
			badge = `<span class="sr-badge sr-badge-pto">PTO</span>`
		case pbool(r.line.Stopped):
			badge = `<span class="sr-badge sr-badge-stp">STP</span>`
		}

		fmt.Fprintf(&sb,
			`<div class="%s" style="%s">`+
				`<span class="sr-pos">%s</span>`+
				`<span class="sr-tla" style="color:%s">%s</span>`+
				`<span class="sr-gap">%s</span>`+
				`%s`+
				`</div>`,
			rowClass, rowStyle,
			html.EscapeString(r.line.Position),
			colour, html.EscapeString(tla),
			gap,
			badge,
		)
	}
	return sb.String()
}
