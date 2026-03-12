package livetiming

import (
	"encoding/json"
	"fmt"
	"time"
)

// FlexTime wraps time.Time to handle F1 date strings that may omit the
// timezone offset (e.g. "2025-07-06T15:00:00" as well as RFC3339).
type FlexTime struct{ time.Time }

var flexTimeFormats = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05.999",
	"2006-01-02T15:04:05",
}

// FlexSlice handles F1 JSON fields that arrive as arrays in the initial state
// snapshot but as index-keyed objects in incremental updates
// (e.g. {"0": {...}, "1": {...}}). Both forms are normalised to map[string]T.
type FlexSlice[T any] map[string]T

func (fs *FlexSlice[T]) UnmarshalJSON(data []byte) error {
	var arr []T
	if err := json.Unmarshal(data, &arr); err == nil {
		*fs = make(FlexSlice[T], len(arr))
		for i, v := range arr {
			(*fs)[fmt.Sprintf("%d", i)] = v
		}
		return nil
	}
	var m map[string]T
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	*fs = m
	return nil
}

func (ft *FlexTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	for _, layout := range flexTimeFormats {
		if t, err := time.Parse(layout, s); err == nil {
			ft.Time = t
			return nil
		}
	}
	return fmt.Errorf("cannot parse time %q", s)
}

type ArchiveStatus struct {
	Status string `json:"Status"`
}

type AudioStream struct {
	Name     string   `json:"Name"`
	Language string   `json:"Language"`
	Uri      string   `json:"Uri"`
	Path     string   `json:"Path"`
	Utc      FlexTime `json:"Utc"`
}

type AudioStreams struct {
	Streams []AudioStream `json:"Streams"`
}

type CarDataChannel struct {
	Rpm      int `json:"0"`
	Speed    int `json:"2"`
	Gear     int `json:"3"`
	Throttle int `json:"4"`
	Brake    int `json:"5"`
	Drs      int `json:"45"`
}

type Car struct {
	Channels CarDataChannel `json:"Channels"`
}

type CarDataEntry struct {
	Utc  FlexTime       `json:"Utc"`
	Cars map[string]Car `json:"Cars"`
}

type CarData struct {
	Entries []CarDataEntry `json:"Entries"`
}

type DriverPrediction struct {
	RacingNumber      string  `json:"RacingNumber"`
	CurrentPosition   int     `json:"CurrentPosition"`
	PredictedPosition int     `json:"PredictedPosition"`
	CurrentPoints     float64 `json:"CurrentPoints"`
	PredictedPoints   float64 `json:"PredictedPoints"`
}

type TeamPrediction struct {
	TeamName          string  `json:"TeamName"`
	CurrentPosition   int     `json:"CurrentPosition"`
	PredictedPosition int     `json:"PredictedPosition"`
	CurrentPoints     float64 `json:"CurrentPoints"`
	PredictedPoints   float64 `json:"PredictedPoints"`
}

type ChampionshipPrediction struct {
	Drivers map[string]DriverPrediction `json:"Drivers"`
	Teams   map[string]TeamPrediction   `json:"Teams"`
}

type ContentStream struct {
	Type     string   `json:"Type"`
	Name     string   `json:"Name"`
	Language string   `json:"Language"`
	Uri      string   `json:"Uri"`
	Utc      FlexTime `json:"Utc"`
	Path     string   `json:"Path,omitempty"`
}

type ContentStreams struct {
	Streams FlexSlice[ContentStream] `json:"Streams"`
}

type Tyre struct {
	Compound string `json:"Compound"`
	New      bool   `json:"New"`
}

type CurrentTyres struct {
	Tyres map[string]Tyre `json:"Tyres"`
}

type Driver struct {
	RacingNumber  string `json:"RacingNumber"`
	BroadcastName string `json:"BroadcastName"`
	FullName      string `json:"FullName"`
	Tla           string `json:"Tla"`
	Line          int    `json:"Line"`
	TeamName      string `json:"TeamName"`
	TeamColour    string `json:"TeamColour"`
	FirstName     string `json:"FirstName"`
	LastName      string `json:"LastName"`
	Reference     string `json:"Reference"`
	HeadshotUrl   string `json:"HeadshotUrl"`
}

type DriverList map[string]Driver

// UnmarshalJSON handles live delta updates where F1 sends `false` instead of
// a Driver object for entries that are being cleared from the list.
func (dl *DriverList) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*dl = make(DriverList, len(raw))
	for k, v := range raw {
		var d Driver
		if err := json.Unmarshal(v, &d); err != nil {
			continue // skip non-object values (e.g. false)
		}
		(*dl)[k] = d
	}
	return nil
}

type DriverRaceInfoEntry struct {
	RacingNumber  string `json:"RacingNumber"`
	Position      string `json:"Position"`
	Gap           string `json:"Gap"`
	Interval      string `json:"Interval"`
	PitStops      int    `json:"PitStops"`
	Catching      int    `json:"Catching"`
	OvertakeState int    `json:"OvertakeState"`
	IsOut         bool   `json:"IsOut"`
}

type DriverRaceInfo map[string]DriverRaceInfoEntry

type DriverScoreKey struct {
	Category string `json:"Category"`
	Name     string `json:"Name"`
}

type DriverScore struct {
	Keys   []DriverScoreKey       `json:"Keys"`
	Scores map[string][][]float64 `json:"Scores"`
}

type ExtrapolatedClock struct {
	Utc           FlexTime `json:"Utc"`
	Remaining     string   `json:"Remaining"`
	Extrapolating bool     `json:"Extrapolating"`
}

type Heartbeat struct {
	Utc FlexTime `json:"Utc"`
}

type Feed struct {
	KeyFramePath string `json:"KeyFramePath"`
	StreamPath   string `json:"StreamPath"`
}
type Index struct {
	Feeds struct {
		SessionInfo            Feed `json:"SessionInfo"`
		ArchiveStatus          Feed `json:"ArchiveStatus"`
		TrackStatus            Feed `json:"TrackStatus"`
		SessionData            Feed `json:"SessionData"`
		ContentStreams         Feed `json:"ContentStreams"`
		ChampionshipPrediction Feed `json:"ChampionshipPrediction"`
		AudioStreams           Feed `json:"AudioStreams"`
		TimingDataF1           Feed `json:"TimingDataF1"`
		TimingData             Feed `json:"TimingData"`
		DriverList             Feed `json:"DriverList"`
		SPFeed                 Feed `json:"SPFeed"`
		LapSeries              Feed `json:"LapSeries"`
		TopThree               Feed `json:"TopThree"`
		TimingAppData          Feed `json:"TimingAppData"`
		TimingStats            Feed `json:"TimingStats"`
		CarDataZ               Feed `json:"CarData.z"`
		PositionZ              Feed `json:"Position.z"`
		ExtrapolatedClock      Feed `json:"ExtrapolatedClock"`
		TyreStintSeries        Feed `json:"TyreStintSeries"`
		DriverRaceInfo         Feed `json:"DriverRaceInfo"`
		LapCount               Feed `json:"LapCount"`
		SessionStatus          Feed `json:"SessionStatus"`
		Heartbeat              Feed `json:"Heartbeat"`
		WeatherData            Feed `json:"WeatherData"`
		WeatherDataSeries      Feed `json:"WeatherDataSeries"`
		TeamRadio              Feed `json:"TeamRadio"`
		TlaRcm                 Feed `json:"TlaRcm"`
		RaceControlMessages    Feed `json:"RaceControlMessages"`
		CurrentTyres           Feed `json:"CurrentTyres"`
		DriverScore            Feed `json:"DriverScore"`
		PitLaneTimeCollection  Feed `json:"PitLaneTimeCollection"`
	} `json:"Feeds"`
}

func (i *Index) String() string {
	buf, _ := json.MarshalIndent(i, "", "  ")
	return string(buf)
}

func (i *Index) GetFeeds() []Feed {
	return []Feed{
		i.Feeds.SessionInfo,
		i.Feeds.ArchiveStatus,
		i.Feeds.TrackStatus,
		i.Feeds.SessionData,
		i.Feeds.ContentStreams,
		i.Feeds.ChampionshipPrediction,
		i.Feeds.AudioStreams,
		i.Feeds.TimingDataF1,
		i.Feeds.TimingData,
		i.Feeds.DriverList,
		i.Feeds.SPFeed,
		i.Feeds.LapSeries,
		i.Feeds.TopThree,
		i.Feeds.TimingAppData,
		i.Feeds.TimingStats,
		i.Feeds.CarDataZ,
		i.Feeds.PositionZ,
		i.Feeds.ExtrapolatedClock,
		i.Feeds.TyreStintSeries,
		i.Feeds.DriverRaceInfo,
		i.Feeds.LapCount,
		i.Feeds.SessionStatus,
		i.Feeds.Heartbeat,
		i.Feeds.WeatherData,
		i.Feeds.WeatherDataSeries,
		i.Feeds.TeamRadio,
		i.Feeds.TlaRcm,
		i.Feeds.RaceControlMessages,
		i.Feeds.CurrentTyres,
		i.Feeds.DriverScore,
		i.Feeds.PitLaneTimeCollection,
	}
}

type LapCount struct {
	CurrentLap int `json:"CurrentLap"`
	TotalLaps  int `json:"TotalLaps"`
}

type LapSeriesEntry struct {
	RacingNumber string            `json:"RacingNumber"`
	LapPosition  FlexSlice[string] `json:"LapPosition"`
}

type LapSeries map[string]LapSeriesEntry

// PitLaneTimeCollection is a message that is streamed, but no data has been observed yet.
type PitLaneTimeCollection struct{}

// SPFeed is a message that is streamed, but no data has been observed yet.
type SPFeed struct{}

type PositionEntry struct {
	Status string `json:"Status"`
	X      int    `json:"X"`
	Y      int    `json:"Y"`
	Z      int    `json:"Z"`
}

type PositionTimestamp struct {
	Timestamp string                   `json:"Timestamp"`
	Entries   map[string]PositionEntry `json:"Entries"`
}

type PositionData struct {
	Utc      FlexTime            `json:"Utc"`
	Position []PositionTimestamp `json:"Position"`
}

type RaceControlMessage struct {
	Utc      FlexTime `json:"Utc"`
	Lap      int      `json:"Lap"`
	Category string   `json:"Category"`
	Message  string   `json:"Message"`
	Flag     string   `json:"Flag,omitempty"`
	Scope    string   `json:"Scope,omitempty"`
	Sector   int      `json:"Sector,omitempty"`
}

type RaceControlMessages struct {
	Messages FlexSlice[RaceControlMessage] `json:"Messages"`
}

type SessionDataPoint struct {
	Utc FlexTime `json:"Utc"`
	Lap int      `json:"Lap"`
}

type SessionData struct {
	Series FlexSlice[SessionDataPoint] `json:"Series"`
}

type Meeting struct {
	Key          int    `json:"Key"`
	Name         string `json:"Name"`
	OfficialName string `json:"OfficialName"`
	Location     string `json:"Location"`
	Number       int    `json:"Number"`
	Country      struct {
		Key  int    `json:"Key"`
		Code string `json:"Code"`
		Name string `json:"Name"`
	} `json:"Country"`
	Circuit struct {
		Key       int    `json:"Key"`
		ShortName string `json:"ShortName"`
	} `json:"Circuit"`
}

type SessionInfo struct {
	Meeting       Meeting       `json:"Meeting"`
	SessionStatus string        `json:"SessionStatus"`
	ArchiveStatus ArchiveStatus `json:"ArchiveStatus"`
	Key           int           `json:"Key"`
	Type          string        `json:"Type"`
	Name          string        `json:"Name"`
	StartDate     FlexTime      `json:"StartDate"`
	EndDate       FlexTime      `json:"EndDate"`
	GmtOffset     string        `json:"GmtOffset"`
	Path          string        `json:"Path"`
}

type SessionStatus struct {
	Status  string `json:"Status"`
	Started string `json:"Started"`
}

type TeamRadioCapture struct {
	Utc          FlexTime `json:"Utc"`
	RacingNumber string   `json:"RacingNumber"`
	Path         string   `json:"Path"`
}

type TeamRadio struct {
	Captures FlexSlice[TeamRadioCapture] `json:"Captures"`
}

type Stint struct {
	LapTime         string `json:"LapTime"`
	LapNumber       int    `json:"LapNumber"`
	LapFlags        int    `json:"LapFlags"`
	Compound        string `json:"Compound"`
	New             string `json:"New"`
	TyresNotChanged string `json:"TyresNotChanged"`
	TotalLaps       int    `json:"TotalLaps"`
	StartLaps       int    `json:"StartLaps"`
}

type TimingAppDataLine struct {
	RacingNumber string           `json:"RacingNumber"`
	Line         int              `json:"Line"`
	GridPos      string           `json:"GridPos"`
	Stints       FlexSlice[Stint] `json:"Stints"`
}

type TimingAppData struct {
	Lines map[string]TimingAppDataLine `json:"Lines"`
}

type Interval struct {
	Value    string `json:"Value"`
	Catching bool   `json:"Catching"`
}

type Segment struct {
	Status int `json:"Status"`
}

type Sector struct {
	Stopped         bool               `json:"Stopped"`
	PreviousValue   string             `json:"PreviousValue"`
	Segments        FlexSlice[Segment] `json:"Segments"`
	Value           string             `json:"Value"`
	Status          int                `json:"Status"`
	OverallFastest  bool               `json:"OverallFastest"`
	PersonalFastest bool               `json:"PersonalFastest"`
}

type Speed struct {
	Value           string `json:"Value"`
	Status          int    `json:"Status"`
	OverallFastest  bool   `json:"OverallFastest"`
	PersonalFastest bool   `json:"PersonalFastest"`
}

type LapTime struct {
	Value string `json:"Value"`
	Lap   int    `json:"Lap"`
}

type TimingDataLine struct {
	GapToLeader             string            `json:"GapToLeader"`
	IntervalToPositionAhead Interval          `json:"IntervalToPositionAhead"`
	Line                    int               `json:"Line"`
	Position                string            `json:"Position"`
	ShowPosition            bool              `json:"ShowPosition"`
	RacingNumber            string            `json:"RacingNumber"`
	Retired                 *bool             `json:"Retired"`
	InPit                   *bool             `json:"InPit"`
	PitOut                  *bool             `json:"PitOut"`
	Stopped                 *bool             `json:"Stopped"`
	Status                  int               `json:"Status"`
	NumberOfLaps            int               `json:"NumberOfLaps"`
	NumberOfPitStops        int               `json:"NumberOfPitStops"`
	Sectors                 FlexSlice[Sector] `json:"Sectors"`
	Speeds                  map[string]Speed  `json:"Speeds"`
	BestLapTime             LapTime           `json:"BestLapTime"`
	LastLapTime             Speed             `json:"LastLapTime"`
}

type TimingData struct {
	Lines map[string]TimingDataLine `json:"Lines"`
}

type Best struct {
	Position int    `json:"Position"`
	Value    string `json:"Value"`
}

type PersonalBestLapTime struct {
	Lap      int    `json:"Lap"`
	Position int    `json:"Position"`
	Value    string `json:"Value"`
}

type TimingStatsLine struct {
	Line                int                 `json:"Line"`
	RacingNumber        string              `json:"RacingNumber"`
	PersonalBestLapTime PersonalBestLapTime `json:"PersonalBestLapTime"`
	BestSectors         FlexSlice[Best]     `json:"BestSectors"`
	BestSpeeds          map[string]Best     `json:"BestSpeeds"`
}

type TimingStats struct {
	Withheld bool                       `json:"Withheld"`
	Lines    map[string]TimingStatsLine `json:"Lines"`
}

type TlaRcm struct {
	Timestamp FlexTime `json:"Timestamp"`
	Message   string   `json:"Message"`
}

type TopThreeLine struct {
	Position        string `json:"Position"`
	ShowPosition    bool   `json:"ShowPosition"`
	RacingNumber    string `json:"RacingNumber"`
	Tla             string `json:"Tla"`
	BroadcastName   string `json:"BroadcastName"`
	FullName        string `json:"FullName"`
	FirstName       string `json:"FirstName"`
	LastName        string `json:"LastName"`
	Reference       string `json:"Reference"`
	Team            string `json:"Team"`
	TeamColour      string `json:"TeamColour"`
	LapTime         string `json:"LapTime"`
	LapState        int    `json:"LapState"`
	DiffToAhead     string `json:"DiffToAhead"`
	DiffToLeader    string `json:"DiffToLeader"`
	OverallFastest  bool   `json:"OverallFastest"`
	PersonalFastest bool   `json:"PersonalFastest"`
}

type TopThree struct {
	Withheld bool                    `json:"Withheld"`
	Lines    FlexSlice[TopThreeLine] `json:"Lines"`
}

type TrackStatus struct {
	Status  string `json:"Status"`
	Message string `json:"Message"`
}

type TyreStint struct {
	Compound        string `json:"Compound"`
	New             string `json:"New"`
	TyresNotChanged string `json:"TyresNotChanged"`
	TotalLaps       int    `json:"TotalLaps"`
	StartLaps       int    `json:"StartLaps"`
}

type TyreStintSeries struct {
	Stints map[string]FlexSlice[TyreStint] `json:"Stints"`
}

type WeatherData struct {
	AirTemp       string `json:"AirTemp"`
	Humidity      string `json:"Humidity"`
	Pressure      string `json:"Pressure"`
	Rainfall      string `json:"Rainfall"`
	TrackTemp     string `json:"TrackTemp"`
	WindDirection string `json:"WindDirection"`
	WindSpeed     string `json:"WindSpeed"`
}

type WeatherDataPoint struct {
	Timestamp FlexTime    `json:"Timestamp"`
	Weather   WeatherData `json:"Weather"`
}

type WeatherDataSeries struct {
	Series FlexSlice[WeatherDataPoint] `json:"Series"`
}

func (h *Heartbeat) String() string {
	return fmt.Sprintf("heartbeat utc=%s", h.Utc.Format("15:04:05"))
}

func (c *CarData) String() string {
	return fmt.Sprintf("car data cars=%d", len(c.Entries))
}

func (p *PositionData) String() string {
	return fmt.Sprintf("position data entries=%d", len(p.Position))
}

func (s *SessionInfo) String() string {
	return fmt.Sprintf("session info meeting=%s session=%s", s.Meeting.Name, s.Name)
}

func (t *TimingData) String() string {
	return fmt.Sprintf("timing data drivers=%d", len(t.Lines))
}

func (t *TopThree) String() string {
	return fmt.Sprintf("top three drivers=%d", len(t.Lines))
}

func (t *TimingStats) String() string {
	return fmt.Sprintf("timing stats drivers=%d", len(t.Lines))
}

func (t *TimingAppData) String() string {
	return fmt.Sprintf("timing app data drivers=%d", len(t.Lines))
}

func (w *WeatherData) String() string {
	return fmt.Sprintf("weather air=%s°C track=%s°C wind=%skm/h rain=%s", w.AirTemp, w.TrackTemp, w.WindSpeed, w.Rainfall)
}

func (t *TrackStatus) String() string {
	return fmt.Sprintf("track status %s %s", t.Status, t.Message)
}

func (d *DriverList) String() string {
	return fmt.Sprintf("driver list drivers=%d", len(*d))
}

func (r *RaceControlMessages) String() string {
	return fmt.Sprintf("race control messages=%d", len(r.Messages))
}

func (s *SessionData) String() string {
	return fmt.Sprintf("session data points=%d", len(s.Series))
}

func (l *LapCount) String() string {
	return fmt.Sprintf("lap count current=%d total=%d", l.CurrentLap, l.TotalLaps)
}

func (e *ExtrapolatedClock) String() string {
	return fmt.Sprintf("clock remaining=%s extrapolating=%v", e.Remaining, e.Extrapolating)
}

func (s *SessionStatus) String() string {
	return fmt.Sprintf("session status %s", s.Status)
}

func (t *TlaRcm) String() string {
	return fmt.Sprintf("tla rcm %s", t.Message)
}

func (t *TeamRadio) String() string {
	if len(t.Captures) > 0 {
		last := t.Captures[fmt.Sprintf("%d", len(t.Captures)-1)]
		return fmt.Sprintf("team radio #%s", last.RacingNumber)
	}
	return "team radio"
}

func (a *ArchiveStatus) String() string {
	return fmt.Sprintf("archive status %s", a.Status)
}

func (a *AudioStreams) String() string {
	return fmt.Sprintf("audio streams streams=%d", len(a.Streams))
}

func (c *ChampionshipPrediction) String() string {
	return fmt.Sprintf("championship prediction drivers=%d teams=%d", len(c.Drivers), len(c.Teams))
}

func (c *ContentStreams) String() string {
	return fmt.Sprintf("content streams streams=%d", len(c.Streams))
}

func (c *CurrentTyres) String() string {
	return fmt.Sprintf("current tyres cars=%d", len(c.Tyres))
}

func (d DriverRaceInfo) String() string {
	return fmt.Sprintf("driver race info drivers=%d", len(d))
}

func (d *DriverScore) String() string {
	return fmt.Sprintf("driver score scores=%d", len(d.Scores))
}

func (l LapSeries) String() string {
	return fmt.Sprintf("lap series drivers=%d", len(l))
}

func (p *PitLaneTimeCollection) String() string {
	return "pit lane time collection"
}

func (s *SPFeed) String() string {
	return "sp feed"
}

func (t *TyreStintSeries) String() string {
	return fmt.Sprintf("tyre stint series drivers=%d", len(t.Stints))
}

func (w *WeatherDataSeries) String() string {
	return fmt.Sprintf("weather data series points=%d", len(w.Series))
}
