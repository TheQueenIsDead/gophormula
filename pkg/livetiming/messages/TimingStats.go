package messages

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
