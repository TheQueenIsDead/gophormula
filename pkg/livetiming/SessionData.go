package livetiming

import "time"

type SessionDataPoint struct {
	Utc time.Time `json:"Utc"`
	Lap int       `json:"Lap"`
}

type SessionData struct {
	Series []SessionDataPoint `json:"Series"`
}
