package livetiming

import "time"

type T struct {
	Streams []struct {
		Name     string    `json:"Name"`
		Language string    `json:"Language"`
		Uri      string    `json:"Uri"`
		Path     string    `json:"Path"`
		Utc      time.Time `json:"Utc"`
	} `json:"Streams"`
}
