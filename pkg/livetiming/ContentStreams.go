package livetiming

type ContentStreams struct {
	Streams []struct {
		Type     string `json:"Type"`
		Name     string `json:"Name"`
		Language string `json:"Language"`
		Uri      string `json:"Uri"`
		Utc      string `json:"Utc"`
		Path     string `json:"Path,omitempty"`
	} `json:"Streams"`
}
