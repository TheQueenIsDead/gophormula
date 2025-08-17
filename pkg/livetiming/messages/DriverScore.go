package messages

type DriverScoreKey struct {
	Category string `json:"Category"`
	Name     string `json:"Name"`
}

type DriverScore struct {
	Keys   []DriverScoreKey       `json:"Keys"`
	Scores map[string][][]float64 `json:"Scores"`
}
