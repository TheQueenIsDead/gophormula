package messages

type TopThreeLine struct {
	Position        string `json:"Position"`
	ShowPosition    bool   `json:"ShowPosition"`
	RacingNumber    string `json:"RacingNumber"`
	Tla             string `json:"Tla"`
	BroadcastName   string `json:"BroadcastName"`
	FullName        string `json:"FullName"`
	FirstName       string `json:"FirstName"`
	LastName        string `json:"LastName"`
	Reference       string `json:"Reference"`
	Team            string `json:"Team"`
	TeamColour      string `json:"TeamColour"`
	LapTime         string `json:"LapTime"`
	LapState        int    `json:"LapState"`
	DiffToAhead     string `json:"DiffToAhead"`
	DiffToLeader    string `json:"DiffToLeader"`
	OverallFastest  bool   `json:"OverallFastest"`
	PersonalFastest bool   `json:"PersonalFastest"`
}

type TopThree struct {
	Withheld bool           `json:"Withheld"`
	Lines    []TopThreeLine `json:"Lines"`
}
