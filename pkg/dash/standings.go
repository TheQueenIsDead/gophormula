package dash

import (
	"fmt"
	"gophormula/pkg/livetiming"
	"gophormula/pkg/session"
	"html"
	"sort"
	"strconv"
	"strings"
)

func pbool(p *bool) bool { return p != nil && *p }

// renderStandings builds the inner HTML for the standings panel from the
// current session state, sorted by race position.
func renderStandings(sess *session.Session) string {
	if len(sess.Standings) == 0 {
		return ""
	}

	type row struct {
		num  string
		line livetiming.TimingDataLine
		pos  int
	}
	rows := make([]row, 0, len(sess.Standings))
	for num, line := range sess.Standings {
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
		driver := sess.Drivers[r.num]
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
