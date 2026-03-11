package dash

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gophormula/pkg/livetiming"
	"gophormula/pkg/replay"
	"html"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// svgW and svgH are the fixed dimensions of the track SVG canvas.
const svgW, svgH = 800, 600

// ReplayHandler returns an HTTP handler that starts a session replay and
// streams updates to all connected SSE clients.
func (h *Hub) ReplayHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Query().Get("path")
		if path == "" {
			http.Error(w, "missing path", http.StatusBadRequest)
			return
		}

		r2 := replay.New()
		if err := r2.ParseGlob(filepath.Join(path, "*")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if seekStr := r.URL.Query().Get("seek"); seekStr != "" {
			if d, err := time.ParseDuration(seekStr); err == nil {
				r2.SeekTo(d)
			}
		}
		bounds := r2.ScanPositionBounds()
		trackSVG := fetchAndBuildTrackSVG(path, bounds)
		ch := r2.StartAndSubscribe()
		go func() {
			standings := newStandingsState()
			status := newStatusState()
			catchingUp := true
			for m := range ch {
				msg := m.(replay.Message)
				// Always accumulate state even during catch-up.
				rerender := false
				switch v := msg.Value.(type) {
				case *livetiming.TimingData:
					if msg.Timestamp != nil {
						standings.mergeTimingData(v)
						rerender = true
					}
				case *livetiming.DriverList:
					standings.mergeDriverList(v)
					rerender = true
				}
				if msg.Timestamp != nil {
					status.merge(msg.Value)
				}
				// During catch-up, skip all UI updates.
				if msg.Catchup {
					continue
				}
				// First real-time message: flush accumulated status and standings.
				if catchingUp {
					catchingUp = false
					status.flush(h)
					if s := standings.render(); s != "" {
						h.send("standings-panel", "inner", s)
					}
				}
				if msg.Timestamp != nil {
					h.BroadcastStatus("status-time", msg.Timestamp.Format("15:04:05"))
				}
				if pd, ok := msg.Value.(*livetiming.PositionData); ok {
					h.BroadcastCars(buildCarsSVG(pd, bounds, standings.drivers, trackSVG))
					continue
				}
				if rerender {
					if s := standings.render(); s != "" {
						h.send("standings-panel", "inner", s)
					}
				}
				if msg.Timestamp != nil {
					updateStatus(h, msg.Value)
				}
				body := formatMessage(msg.Value)
				if body == "" {
					continue
				}
				ts := "--:--:--"
				if msg.Timestamp != nil {
					ts = msg.Timestamp.Format("15:04:05")
				}
				h.Broadcast(ts, body)
			}
		}()

		SessionStarted(w, r, filepath.Base(path))
	}
}

// fetchAndBuildTrackSVG reads SessionInfo.json from the session directory,
// fetches the circuit map from the Multiviewer API, and returns an SVG polyline
// fragment for the track outline. Falls back to an empty string on any error.
func fetchAndBuildTrackSVG(sessionPath string, b replay.PositionBounds) string {
	raw, err := os.ReadFile(filepath.Join(sessionPath, "SessionInfo.json"))
	if err != nil {
		log.Printf("circuit map: reading SessionInfo.json: %v", err)
		return ""
	}
	raw = bytes.TrimPrefix(raw, []byte{0xEF, 0xBB, 0xBF}) // strip UTF-8 BOM
	var si livetiming.SessionInfo
	if err := json.Unmarshal(raw, &si); err != nil {
		log.Printf("circuit map: parsing SessionInfo: %v", err)
		return ""
	}
	year := si.StartDate.Year()
	if year == 0 {
		year = time.Now().Year()
	}
	cm, err := livetiming.FetchCircuitMap(si.Meeting.Circuit.Key, year)
	if err != nil {
		log.Printf("circuit map: fetch failed: %v", err)
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

// buildCarsSVG converts a PositionData snapshot into a full SVG element with
// car positions normalised to the fixed svgW×svgH canvas using the
// pre-scanned circuit bounds. drivers provides team colours and TLAs.
// trackSVG is the pre-built circuit outline fragment from buildTrackSVG.
func buildCarsSVG(pd *livetiming.PositionData, b replay.PositionBounds, drivers map[string]livetiming.Driver, trackSVG string) string {
	if len(pd.Position) == 0 || !b.Valid {
		return ""
	}
	last := pd.Position[len(pd.Position)-1]
	rangeX := b.MaxX - b.MinX
	rangeY := b.MaxY - b.MinY
	if rangeX == 0 {
		rangeX = 1
	}
	if rangeY == 0 {
		rangeY = 1
	}

	const pad = 40 // pixel padding inside the canvas
	usableW := svgW - 2*pad
	usableH := svgH - 2*pad

	var cars strings.Builder
	for num, e := range last.Entries {
		// normalise: SVG Y is flipped relative to F1 Y
		sx := float64(pad) + float64(e.X-b.MinX)/float64(rangeX)*float64(usableW)
		sy := float64(pad) + (1.0-float64(e.Y-b.MinY)/float64(rangeY))*float64(usableH)
		dotColor := "#ffffff"
		if e.Status == "OffTrack" {
			dotColor = "#555555"
		} else if d, ok := drivers[num]; ok && d.TeamColour != "" {
			dotColor = "#" + d.TeamColour
		}
		label := num
		if d, ok := drivers[num]; ok && d.Tla != "" {
			label = d.Tla
		}
		textColor := "#111111"
		if e.Status == "OffTrack" {
			textColor = "#999999"
		}
		fmt.Fprintf(&cars,
			`<circle cx="%.1f" cy="%.1f" r="9" fill="%s" stroke="#111" stroke-width="1"></circle>`+
				`<text x="%.1f" y="%.1f" text-anchor="middle" font-size="7" fill="%s" font-family="monospace" font-weight="bold">%s</text>`,
			sx, sy, dotColor, sx, sy+2.5, textColor, html.EscapeString(label))
	}
	return fmt.Sprintf(
		`<svg id="track-plot" viewBox="0 0 %d %d" preserveAspectRatio="xMidYMid meet">`+
			`<rect width="%d" height="%d" fill="#0a0a0a"></rect>`+
			`<g id="track">%s</g>`+
			`<g id="cars">%s</g>`+
			`</svg>`,
		svgW, svgH, svgW, svgH, trackSVG, cars.String(),
	)
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
	switch v := msg.(type) {
	case *livetiming.Heartbeat:
		return fmt.Sprintf("Heartbeat  utc=%s", v.Utc.Format("15:04:05"))
	case *livetiming.LapCount:
		return fmt.Sprintf("LapCount  lap=%d/%d", v.CurrentLap, v.TotalLaps)
	case *livetiming.TrackStatus:
		return fmt.Sprintf("TrackStatus  %s  %s", v.Status, v.Message)
	case *livetiming.SessionStatus:
		return fmt.Sprintf("SessionStatus  %s", v.Status)
	case *livetiming.SessionInfo:
		return fmt.Sprintf("SessionInfo  %s — %s", v.Meeting.Name, v.Name)
	case *livetiming.ExtrapolatedClock:
		return fmt.Sprintf("Clock  remaining=%s  extrapolating=%v", v.Remaining, v.Extrapolating)
	case *livetiming.WeatherData:
		return fmt.Sprintf("WeatherData  air=%s°C  track=%s°C  wind=%skm/h  rain=%s", v.AirTemp, v.TrackTemp, v.WindSpeed, v.Rainfall)
	case *livetiming.RaceControlMessages:
		if len(v.Messages) > 0 {
			last := v.Messages[fmt.Sprintf("%d", len(v.Messages)-1)]
			return fmt.Sprintf("RaceControl  [%s] %s", last.Category, last.Message)
		}
		return "RaceControl"
	case *livetiming.TlaRcm:
		return fmt.Sprintf("TlaRcm  %s", v.Message)
	case *livetiming.TeamRadio:
		if len(v.Captures) > 0 {
			last := v.Captures[fmt.Sprintf("%d", len(v.Captures)-1)]
			return fmt.Sprintf("TeamRadio  #%s", last.RacingNumber)
		}
		return "TeamRadio"
	case *livetiming.CarData:
		if len(v.Entries) == 0 {
			return ""
		}
		last := v.Entries[len(v.Entries)-1]
		parts := make([]string, 0, len(last.Cars))
		for num, car := range last.Cars {
			ch := car.Channels
			parts = append(parts, fmt.Sprintf("#%s %dkm/h G%d T%d%% B%d%% DRS%d",
				num, ch.Speed, ch.Gear, ch.Throttle, ch.Brake, ch.Drs))
		}
		return "CarData  " + strings.Join(parts, "  |  ")
	case *livetiming.DriverList:
		return fmt.Sprintf("DriverList  %d drivers", len(*v))
	case *livetiming.TimingData:
		if len(v.Lines) > 3 {
			return fmt.Sprintf("TimingData  %d drivers", len(v.Lines))
		}
		driverParts := make([]string, 0, len(v.Lines))
		for num, line := range v.Lines {
			fields := []string{fmt.Sprintf("#%s", num)}
			if line.Position != "" {
				fields = append(fields, fmt.Sprintf("P%s", line.Position))
			}
			if line.NumberOfLaps > 0 {
				fields = append(fields, fmt.Sprintf("L%d", line.NumberOfLaps))
			}
			if line.LastLapTime.Value != "" {
				fields = append(fields, fmt.Sprintf("last=%s", line.LastLapTime.Value))
			}
			if line.BestLapTime.Value != "" {
				fields = append(fields, fmt.Sprintf("best=%s", line.BestLapTime.Value))
			}
			if line.GapToLeader != "" {
				fields = append(fields, fmt.Sprintf("gap=%s", line.GapToLeader))
			}
			if line.IntervalToPositionAhead.Value != "" {
				fields = append(fields, fmt.Sprintf("int=%s", line.IntervalToPositionAhead.Value))
			}
			if line.NumberOfPitStops > 0 {
				fields = append(fields, fmt.Sprintf("pits=%d", line.NumberOfPitStops))
			}
			sectorVals := make([]string, 0, len(line.Sectors))
			for _, s := range line.Sectors {
				if s.Value != "" {
					sectorVals = append(sectorVals, s.Value)
				}
			}
			if len(sectorVals) > 0 {
				fields = append(fields, "S:"+strings.Join(sectorVals, "/"))
			}
			if pbool(line.Retired) {
				fields = append(fields, "[OUT]")
			} else if pbool(line.InPit) {
				fields = append(fields, "[PIT]")
			} else if pbool(line.PitOut) {
				fields = append(fields, "[PIT OUT]")
			} else if pbool(line.Stopped) {
				fields = append(fields, "[STOPPED]")
			}
			if len(fields) > 1 { // skip empty init pings that carry only the racing number
				driverParts = append(driverParts, strings.Join(fields, " "))
			}
		}
		if len(driverParts) == 0 {
			return ""
		}
		return "TimingData  " + strings.Join(driverParts, "  |  ")
	case *livetiming.TimingAppData:
		if len(v.Lines) == 1 {
			for num, line := range v.Lines {
				if len(line.Stints) > 0 {
					last := line.Stints[fmt.Sprintf("%d", len(line.Stints)-1)]
					return fmt.Sprintf("TimingAppData  #%s  %s  laps=%d", num, last.Compound, last.TotalLaps)
				}
				return fmt.Sprintf("TimingAppData  #%s", num)
			}
		}
		return fmt.Sprintf("TimingAppData  %d drivers", len(v.Lines))
	case *livetiming.TimingStats:
		if len(v.Lines) == 1 {
			for num, line := range v.Lines {
				return fmt.Sprintf("TimingStats  #%s  best=%s", num, line.PersonalBestLapTime.Value)
			}
		}
		return fmt.Sprintf("TimingStats  %d drivers", len(v.Lines))
	case *livetiming.TopThree:
		parts := make([]string, 0, 3)
		for i := range 3 {
			if line, ok := v.Lines[fmt.Sprintf("%d", i)]; ok && line.Tla != "" {
				parts = append(parts, fmt.Sprintf("P%s %s %s", line.Position, line.Tla, line.LapTime))
			}
		}
		if len(parts) > 0 {
			return "TopThree  " + strings.Join(parts, "  |  ")
		}
		return "TopThree"
	default:
		return fmt.Sprintf("%T", msg)
	}
}
