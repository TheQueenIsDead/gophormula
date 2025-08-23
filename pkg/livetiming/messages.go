package livetiming

import (
	"encoding/json"
	"time"
)

type ArchiveStatus struct {
	Status string `json:"Status"`
}

type AudioStream struct {
	Name     string    `json:"Name"`
	Language string    `json:"Language"`
	Uri      string    `json:"Uri"`
	Path     string    `json:"Path"`
	Utc      time.Time `json:"Utc"`
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
	Utc  time.Time      `json:"Utc"`
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

type ContentStreams struct {
	Streams []struct {
		Type     string `json:"Type"`
		Name     string `json:"Name"`
		Language string `json:"Language"`
		Uri      string `json:"Uri"`
		Utc      string `json:"Utc"`
		Path     string `json:"Path,omitempty"`
	} `json:"Streams"`
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
	Utc           time.Time `json:"Utc"`
	Remaining     string    `json:"Remaining"`
	Extrapolating bool      `json:"Extrapolating"`
}

type Heartbeat struct {
	Utc time.Time `json:"Utc"`
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
	RacingNumber string   `json:"RacingNumber"`
	LapPosition  []string `json:"LapPosition"`
}

type LapSeries map[string]LapSeriesEntry

// PitLaneTimeCollection is a message that is streamed, but no data has been observed yet.
type PitLaneTimeCollection struct{}

// PositionData defines the structure for position data.
// The actual data is compressed and needs to be handled accordingly.
type PositionData struct {
	// Placeholder for the compressed data structure
}

type RaceControlMessage struct {
	Utc      time.Time `json:"Utc"`
	Lap      int       `json:"Lap"`
	Category string    `json:"Category"`
	Message  string    `json:"Message"`
	Flag     string    `json:"Flag,omitempty"`
	Scope    string    `json:"Scope,omitempty"`
	Sector   int       `json:"Sector,omitempty"`
}

type RaceControlMessages struct {
	Messages []RaceControlMessage `json:"Messages"`
}

type SessionDataPoint struct {
	Utc time.Time `json:"Utc"`
	Lap int       `json:"Lap"`
}

type SessionData struct {
	Series []SessionDataPoint `json:"Series"`
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
	Meeting       Meeting `json:"Meeting"`
	SessionStatus string  `json:"SessionStatus"`
	ArchiveStatus struct {
		Status string `json:"Status"`
	} `json:"ArchiveStatus"`
	Key       int       `json:"Key"`
	Type      string    `json:"Type"`
	Name      string    `json:"Name"`
	StartDate time.Time `json:"StartDate"`
	EndDate   time.Time `json:"EndDate"`
	GmtOffset string    `json:"GmtOffset"`
	Path      string    `json:"Path"`
}

type SessionStatus struct {
	Status  string `json:"Status"`
	Started string `json:"Started"`
}

type TeamRadioCapture struct {
	Utc          time.Time `json:"Utc"`
	RacingNumber string    `json:"RacingNumber"`
	Path         string    `json:"Path"`
}

type TeamRadio struct {
	Captures []TeamRadioCapture `json:"Captures"`
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
	RacingNumber string  `json:"RacingNumber"`
	Line         int     `json:"Line"`
	GridPos      string  `json:"GridPos"`
	Stints       []Stint `json:"Stints"`
}

type TimingAppData struct {
	Lines map[string]TimingAppDataLine `json:"Lines"`
}

type Interval struct {
	Value    string `json:"Value"`
	Catching bool   `json:"Catching"`
}

type Sector struct {
	Stopped       bool   `json:"Stopped"`
	PreviousValue string `json:"PreviousValue"`
	Segments      []struct {
		Status int `json:"Status"`
	} `json:"Segments"`
	Value           string `json:"Value"`
	Status          int    `json:"Status"`
	OverallFastest  bool   `json:"OverallFastest"`
	PersonalFastest bool   `json:"PersonalFastest"`
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
	GapToLeader             string           `json:"GapToLeader"`
	IntervalToPositionAhead Interval         `json:"IntervalToPositionAhead"`
	Line                    int              `json:"Line"`
	Position                string           `json:"Position"`
	ShowPosition            bool             `json:"ShowPosition"`
	RacingNumber            string           `json:"RacingNumber"`
	Retired                 bool             `json:"Retired"`
	InPit                   bool             `json:"InPit"`
	PitOut                  bool             `json:"PitOut"`
	Stopped                 bool             `json:"Stopped"`
	Status                  int              `json:"Status"`
	NumberOfLaps            int              `json:"NumberOfLaps"`
	NumberOfPitStops        int              `json:"NumberOfPitStops"`
	Sectors                 []Sector         `json:"Sectors"`
	Speeds                  map[string]Speed `json:"Speeds"`
	BestLapTime             LapTime          `json:"BestLapTime"`
	LastLapTime             struct {
		Value           string `json:"Value"`
		Status          int    `json:"Status"`
		OverallFastest  bool   `json:"OverallFastest"`
		PersonalFastest bool   `json:"PersonalFastest"`
	} `json:"LastLapTime"`
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
	BestSectors         []Best              `json:"BestSectors"`
	BestSpeeds          map[string]Best     `json:"BestSpeeds"`
}

type TimingStats struct {
	Withheld bool                       `json:"Withheld"`
	Lines    map[string]TimingStatsLine `json:"Lines"`
}

type TlaRcm struct {
	Timestamp time.Time `json:"Timestamp"`
	Message   string    `json:"Message"`
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
	Withheld bool           `json:"Withheld"`
	Lines    []TopThreeLine `json:"Lines"`
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
	Stints map[string][]TyreStint `json:"Stints"`
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
	Timestamp time.Time   `json:"Timestamp"`
	Weather   WeatherData `json:"Weather"`
}

type WeatherDataSeries struct {
	Series []WeatherDataPoint `json:"Series"`
}
