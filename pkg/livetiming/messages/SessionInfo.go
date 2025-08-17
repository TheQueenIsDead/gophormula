package messages

import "time"

type Meeting struct {
	Key          int    `json:"Key"`
	Name         string `json:"Name"`
	OfficialName string `json:"OfficialName"`
	Location     string `json:"Location"`
	Number       int    `json:"Number"`
	Country      struct {
		Key  int    `json:"Key"`
		Code string `json:"Code"`
		Name string `json:"Name"`
	} `json:"Country"`
	Circuit struct {
		Key       int    `json:"Key"`
		ShortName string `json:"ShortName"`
	} `json:"Circuit"`
}

type SessionInfo struct {
	Meeting       Meeting `json:"Meeting"`
	SessionStatus string  `json:"SessionStatus"`
	ArchiveStatus struct {
		Status string `json:"Status"`
	} `json:"ArchiveStatus"`
	Key       int       `json:"Key"`
	Type      string    `json:"Type"`
	Name      string    `json:"Name"`
	StartDate time.Time `json:"StartDate"`
	EndDate   time.Time `json:"EndDate"`
	GmtOffset string    `json:"GmtOffset"`
	Path      string    `json:"Path"`
}
