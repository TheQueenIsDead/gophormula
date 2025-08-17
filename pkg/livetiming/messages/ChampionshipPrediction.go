package messages

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
