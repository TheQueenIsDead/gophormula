package livetiming

import "time"

type CarDataChannel struct {
	Rpm      int `json:"0"`
	Speed    int `json:"2"`
	Gear     int `json:"3"`
	Throttle int `json:"4"`
	Brake    int `json:"5"`
	Drs      int `json:"45"`
}

type Car struct {
	Channels CarDataChannel `json:"Channels"`
}

type CarDataEntry struct {
	Utc  time.Time      `json:"Utc"`
	Cars map[string]Car `json:"Cars"`
}

type CarData struct {
	Entries []CarDataEntry `json:"Entries"`
}
