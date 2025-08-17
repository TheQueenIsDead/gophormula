package messages

import "time"

type Heartbeat struct {
	Utc time.Time `json:"Utc"`
}
