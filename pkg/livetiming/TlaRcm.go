package livetiming

import "time"

type TlaRcm struct {
	Timestamp time.Time `json:"Timestamp"`
	Message   string    `json:"Message"`
}
