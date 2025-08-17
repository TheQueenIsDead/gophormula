package livetiming

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
