package frontend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gophormula/pkg/livetiming"
	"gophormula/pkg/replay"
	"html"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// carDot is the JSON payload sent to the client for a single car.
type carDot struct {
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Color string  `json:"color"`
	Label string  `json:"label"`
	Off   bool    `json:"off"`
}

// svgW and svgH are the fixed dimensions of the track SVG canvas.
const svgW, svgH = 800, 600

// fetchAndBuildTrackSVG reads SessionInfo.json from the session directory,
// fetches the circuit map from the Multiviewer API, and returns an SVG polyline
// fragment for the track outline. Falls back to an empty string on any error.
func fetchAndBuildTrackSVG(sessionPath string, b replay.PositionBounds) string {
	raw, err := os.ReadFile(filepath.Join(sessionPath, "SessionInfo.json"))
	if err != nil {
		slog.Warn("circuit map: reading SessionInfo.json", "err", err)
		return ""
	}
	raw = bytes.TrimPrefix(raw, []byte{0xEF, 0xBB, 0xBF}) // strip UTF-8 BOM
	var si livetiming.SessionInfo
	if err := json.Unmarshal(raw, &si); err != nil {
		slog.Warn("circuit map: parsing SessionInfo", "err", err)
		return ""
	}
	year := si.StartDate.Year()
	if year == 0 {
		year = time.Now().Year()
	}
	cm, err := livetiming.FetchCircuitMap(si.Meeting.Circuit.Key, year)
	if err != nil {
		slog.Warn("circuit map: fetch failed", "err", err)
		return ""
	}
	return buildTrackSVGFromMap(cm, b)
}

// buildTrackSVGFromMap normalises the Multiviewer circuit map coordinates onto
// the fixed SVG canvas and returns a single polyline element.
func buildTrackSVGFromMap(cm *livetiming.CircuitMap, b replay.PositionBounds) string {
	if cm == nil || len(cm.X) == 0 || !b.Valid {
		return ""
	}
	rangeX := float64(b.MaxX - b.MinX)
	rangeY := float64(b.MaxY - b.MinY)
	if rangeX == 0 {
		rangeX = 1
	}
	if rangeY == 0 {
		rangeY = 1
	}
	const pad = 40
	usableW := float64(svgW - 2*pad)
	usableH := float64(svgH - 2*pad)
	var pts strings.Builder
	for i := range cm.X {
		tx := float64(pad) + (cm.X[i]-float64(b.MinX))/rangeX*usableW
		ty := float64(pad) + (1.0-(cm.Y[i]-float64(b.MinY))/rangeY)*usableH
		if i > 0 {
			pts.WriteByte(' ')
		}
		fmt.Fprintf(&pts, "%.1f,%.1f", tx, ty)
	}
	return fmt.Sprintf(
		`<polyline points="%s" fill="none" stroke="#333" stroke-width="8" stroke-linejoin="round" stroke-linecap="round"></polyline>`+
			`<polyline points="%s" fill="none" stroke="#1a1a1a" stroke-width="4" stroke-linejoin="round" stroke-linecap="round"></polyline>`,
		pts.String(), pts.String(),
	)
}

// buildCarsScript converts a single PositionTimestamp frame into a JavaScript
// call to updateCarTargets({...}) with normalised SVG canvas coordinates.
// The client-side rAF loop interpolates smoothly toward the new targets.
func buildCarsScript(frame livetiming.PositionTimestamp, b replay.PositionBounds, drivers map[string]livetiming.Driver) string {
	if !b.Valid {
		return ""
	}
	rangeX := b.MaxX - b.MinX
	rangeY := b.MaxY - b.MinY
	if rangeX == 0 {
		rangeX = 1
	}
	if rangeY == 0 {
		rangeY = 1
	}

	const pad = 40
	usableW := svgW - 2*pad
	usableH := svgH - 2*pad

	cars := make(map[string]carDot, len(frame.Entries))
	for num, e := range frame.Entries {
		sx := float64(pad) + float64(e.X-b.MinX)/float64(rangeX)*float64(usableW)
		sy := float64(pad) + (1.0-float64(e.Y-b.MinY)/float64(rangeY))*float64(usableH)
		color := "#ffffff"
		if d, ok := drivers[num]; ok && d.TeamColour != "" {
			color = "#" + d.TeamColour
		}
		label := num
		if d, ok := drivers[num]; ok && d.Tla != "" {
			label = d.Tla
		}
		cars[num] = carDot{X: sx, Y: sy, Color: color, Label: html.EscapeString(label), Off: e.Status == "OffTrack"}
	}
	data, err := json.Marshal(cars)
	if err != nil {
		return ""
	}
	return "updateCarTargets(" + string(data) + ")"
}

// updateStatus pushes live values to the status bar elements for the message
// types that are displayed there.
func updateStatus(h *Hub, msg any) {
	switch v := msg.(type) {
	case *livetiming.SessionStatus:
		color := sessionStatusColor(v.Status)
		h.BroadcastStatus("status-session",
			fmt.Sprintf(`<span class="status-dot" style="background:%s" data-tip="%s"></span>`,
				color, html.EscapeString(v.Status)))

	case *livetiming.ExtrapolatedClock:
		h.BroadcastStatus("status-remaining", html.EscapeString(v.Remaining))
	case *livetiming.LapCount:
		h.BroadcastStatus("status-lap", fmt.Sprintf("Lap %d / %d", v.CurrentLap, v.TotalLaps))
	case *livetiming.WeatherData:
		h.BroadcastStatus("status-weather",
			fmt.Sprintf("Air %s°C · Track %s°C · Wind %skm/h · Rain %s", v.AirTemp, v.TrackTemp, v.WindSpeed, v.Rainfall))
	case *livetiming.TrackStatus:
		color := trackStatusColor(v.Status)
		h.BroadcastStatus("status-track",
			fmt.Sprintf(`<span style="color:%s;font-weight:bold">%s</span>`, color, html.EscapeString(v.Message)))
	}
}

// sessionStatusColor maps F1 session status strings to dot colours.
func sessionStatusColor(status string) string {
	switch status {
	case "Started":
		return "#00d2be"
	case "Finished", "Finalised", "Ends":
		return "#e10600"
	default: // Inactive, unknown
		return "#ff8c00"
	}
}

// trackStatusColor maps F1 track status codes to display colours.
// 1 = AllClear, 2 = Yellow, 3 = SCDeployed, 4 = SCStopped, 5 = RedFlag,
// 6 = VSCDeployed, 7 = VSCEnding.
func trackStatusColor(status string) string {
	switch status {
	case "1":
		return "#00d2be" // green
	case "2":
		return "#ffd700" // yellow flag
	case "3", "4", "6", "7":
		return "#ffd700" // safety car / VSC
	case "5":
		return "#e10600" // red flag
	default:
		return "#888888"
	}
}

func formatMessage(msg any) string {
	if s, ok := msg.(fmt.Stringer); ok {
		return s.String()
	}
	return fmt.Sprintf("%v", msg)
}

// boundsFromCircuitMap derives PositionBounds from a Multiviewer circuit map.
// The circuit map and F1 position data share the same coordinate space, so the
// circuit extent can normalise live car positions onto the SVG canvas.
func boundsFromCircuitMap(cm *livetiming.CircuitMap) replay.PositionBounds {
	if cm == nil || len(cm.X) == 0 {
		return replay.PositionBounds{}
	}
	minX, maxX := math.MaxInt, math.MinInt
	minY, maxY := math.MaxInt, math.MinInt
	for i := range cm.X {
		x, y := int(math.Round(cm.X[i])), int(math.Round(cm.Y[i]))
		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
	}
	return replay.PositionBounds{MinX: minX, MaxX: maxX, MinY: minY, MaxY: maxY, Valid: true}
}
