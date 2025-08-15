package livetiming

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
