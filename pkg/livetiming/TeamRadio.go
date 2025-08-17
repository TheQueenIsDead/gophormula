package livetiming

import "time"

type TeamRadioCapture struct {
	Utc          time.Time `json:"Utc"`
	RacingNumber string    `json:"RacingNumber"`
	Path         string    `json:"Path"`
}

type TeamRadio struct {
	Captures []TeamRadioCapture `json:"Captures"`
}
