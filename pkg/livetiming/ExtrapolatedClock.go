package livetiming

import "time"

type ExtrapolatedClock struct {
	Utc           time.Time `json:"Utc"`
	Remaining     string    `json:"Remaining"`
	Extrapolating bool      `json:"Extrapolating"`
}
