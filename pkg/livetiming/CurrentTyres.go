package livetiming

type Tyre struct {
	Compound string `json:"Compound"`
	New      bool   `json:"New"`
}

type CurrentTyres struct {
	Tyres map[string]Tyre `json:"Tyres"`
}
