package messages

type TyreStint struct {
	Compound        string `json:"Compound"`
	New             string `json:"New"`
	TyresNotChanged string `json:"TyresNotChanged"`
	TotalLaps       int    `json:"TotalLaps"`
	StartLaps       int    `json:"StartLaps"`
}

type TyreStintSeries struct {
	Stints map[string][]TyreStint `json:"Stints"`
}
