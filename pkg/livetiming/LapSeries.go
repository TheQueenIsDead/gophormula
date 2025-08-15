package livetiming

type LapSeriesEntry struct {
	RacingNumber string   `json:"RacingNumber"`
	LapPosition  []string `json:"LapPosition"`
}

type LapSeries map[string]LapSeriesEntry
