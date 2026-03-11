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

// mergeTimingData applies a delta TimingData update to the accumulated state.
// Only non-zero fields in the delta overwrite existing values; the exception is
// PitOut which clears InPit (car has left the pits).
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
		if delta.Retired {
			existing.Retired = true
		}
		if delta.PitOut {
			existing.InPit = false
			existing.PitOut = true
		} else if delta.InPit {
			existing.InPit = true
			existing.PitOut = false
		}
		if delta.Stopped {
			existing.Stopped = true
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
		ri, rj := rows[i].line.Retired, rows[j].line.Retired
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

		rowStyle := fmt.Sprintf("border-left:3px solid %s", colour)
		rowClass := "sr"
		if r.line.Retired {
			rowClass += " sr-out"
		}

		badge := ""
		switch {
		case r.line.Retired:
			badge = `<span class="sr-badge sr-badge-out">OUT</span>`
		case r.line.InPit:
			badge = `<span class="sr-badge sr-badge-pit">PIT</span>`
		case r.line.PitOut:
			badge = `<span class="sr-badge sr-badge-pto">PTO</span>`
		case r.line.Stopped:
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
