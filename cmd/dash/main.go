package main

import (
	"fmt"
	"gophormula/pkg/dash"
	"gophormula/pkg/livetiming"
	"gophormula/pkg/replay"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// svgW and svgH are the fixed dimensions of the track SVG canvas.
const svgW, svgH = 800, 600

func main() {
	dataDir := "data"
	if len(os.Args) >= 2 {
		dataDir = os.Args[1]
	}

	hub := dash.NewHub(dataDir)

	mux := http.NewServeMux()
	mux.HandleFunc("/", hub.Index)
	mux.HandleFunc("/events", hub.Events)
	mux.HandleFunc("/replay", replayHandler(hub))

	log.Println("listening on :1234")
	if err := http.ListenAndServe(":1234", mux); err != nil {
		panic(err)
	}
}

func replayHandler(hub *dash.Hub) http.HandlerFunc {
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
		bounds := r2.ScanPositionBounds()
		ch := r2.StartAndSubscribe()
		go func() {
			for m := range ch {
				msg := m.(replay.Message)
				if pd, ok := msg.Value.(*livetiming.PositionData); ok {
					hub.BroadcastCars(buildCarsSVG(pd, bounds))
					continue
				}
				body := format(msg.Value)
				if body == "" {
					continue
				}
				ts := "--:--:--"
				if msg.Timestamp != nil {
					ts = msg.Timestamp.Format("15:04:05")
				}
				hub.Broadcast(ts, body)
			}
		}()

		dash.SessionStarted(w, r, filepath.Base(path))
	}
}

// buildCarsSVG converts a PositionData snapshot into SVG circle elements for
// all cars, normalised to the fixed svgW×svgH canvas using the pre-scanned
// circuit bounds.
func buildCarsSVG(pd *livetiming.PositionData, b replay.PositionBounds) string {
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
		color := "#ffffff"
		if e.Status == "OffTrack" {
			color = "#888888"
		}
		fmt.Fprintf(&cars,
			`<circle cx="%.1f" cy="%.1f" r="6" fill="%s" stroke="#111" stroke-width="1"></circle>`+
				`<text x="%.1f" y="%.1f" text-anchor="middle" font-size="8" fill="#111" font-family="monospace">%s</text>`,
			sx, sy, color, sx, sy+3.5, num)
	}
	return fmt.Sprintf(
		`<svg id="track-plot" viewBox="0 0 %d %d" preserveAspectRatio="xMidYMid meet">`+
			`<rect width="%d" height="%d" fill="#0a0a0a"></rect>`+
			`<g id="cars">%s</g>`+
			`</svg>`,
		svgW, svgH, svgW, svgH, cars.String(),
	)
}

func format(msg any) string {
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
			if line.Retired {
				fields = append(fields, "[OUT]")
			} else if line.InPit {
				fields = append(fields, "[PIT]")
			} else if line.PitOut {
				fields = append(fields, "[PIT OUT]")
			} else if line.Stopped {
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
