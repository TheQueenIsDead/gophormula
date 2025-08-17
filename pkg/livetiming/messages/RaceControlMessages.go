package messages

import "time"

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
