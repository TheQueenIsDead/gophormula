package livetiming

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
