package messages

import "time"

type AudioStream struct {
	Name     string    `json:"Name"`
	Language string    `json:"Language"`
	Uri      string    `json:"Uri"`
	Path     string    `json:"Path"`
	Utc      time.Time `json:"Utc"`
}

type AudioStreams struct {
	Streams []AudioStream `json:"Streams"`
}
