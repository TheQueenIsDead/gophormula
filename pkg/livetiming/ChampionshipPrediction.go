package livetiming

type ChampionshipPrediction struct {
	Drivers struct {
		Field1 struct {
			RacingNumber      string  `json:"RacingNumber"`
			CurrentPosition   int     `json:"CurrentPosition"`
			PredictedPosition int     `json:"PredictedPosition"`
			CurrentPoints     float64 `json:"CurrentPoints"`
			PredictedPoints   float64 `json:"PredictedPoints"`
		} `json:"44"`
		Field2 struct {
			RacingNumber      string  `json:"RacingNumber"`
			CurrentPosition   int     `json:"CurrentPosition"`
			PredictedPosition int     `json:"PredictedPosition"`
			CurrentPoints     float64 `json:"CurrentPoints"`
			PredictedPoints   float64 `json:"PredictedPoints"`
		} `json:"33"`
		Field3 struct {
			RacingNumber      string  `json:"RacingNumber"`
			CurrentPosition   int     `json:"CurrentPosition"`
			PredictedPosition int     `json:"PredictedPosition"`
			CurrentPoints     float64 `json:"CurrentPoints"`
			PredictedPoints   float64 `json:"PredictedPoints"`
		} `json:"77"`
		Field4 struct {
			RacingNumber      string  `json:"RacingNumber"`
			CurrentPosition   int     `json:"CurrentPosition"`
			PredictedPosition int     `json:"PredictedPosition"`
			CurrentPoints     float64 `json:"CurrentPoints"`
			PredictedPoints   float64 `json:"PredictedPoints"`
		} `json:"4"`
		Field5 struct {
			RacingNumber      string  `json:"RacingNumber"`
			CurrentPosition   int     `json:"CurrentPosition"`
			PredictedPosition int     `json:"PredictedPosition"`
			CurrentPoints     float64 `json:"CurrentPoints"`
			PredictedPoints   float64 `json:"PredictedPoints"`
		} `json:"11"`
		Field6 struct {
			RacingNumber      string  `json:"RacingNumber"`
			CurrentPosition   int     `json:"CurrentPosition"`
			PredictedPosition int     `json:"PredictedPosition"`
			CurrentPoints     float64 `json:"CurrentPoints"`
			PredictedPoints   float64 `json:"PredictedPoints"`
		} `json:"16"`
		Field7 struct {
			RacingNumber      string  `json:"RacingNumber"`
			CurrentPosition   int     `json:"CurrentPosition"`
			PredictedPosition int     `json:"PredictedPosition"`
			CurrentPoints     float64 `json:"CurrentPoints"`
			PredictedPoints   float64 `json:"PredictedPoints"`
		} `json:"3"`
		Field8 struct {
			RacingNumber      string  `json:"RacingNumber"`
			CurrentPosition   int     `json:"CurrentPosition"`
			PredictedPosition int     `json:"PredictedPosition"`
			CurrentPoints     float64 `json:"CurrentPoints"`
			PredictedPoints   float64 `json:"PredictedPoints"`
		} `json:"55"`
		Field9 struct {
			RacingNumber      string  `json:"RacingNumber"`
			CurrentPosition   int     `json:"CurrentPosition"`
			PredictedPosition int     `json:"PredictedPosition"`
			CurrentPoints     float64 `json:"CurrentPoints"`
			PredictedPoints   float64 `json:"PredictedPoints"`
		} `json:"22"`
		Field10 struct {
			RacingNumber      string  `json:"RacingNumber"`
			CurrentPosition   int     `json:"CurrentPosition"`
			PredictedPosition int     `json:"PredictedPosition"`
			CurrentPoints     float64 `json:"CurrentPoints"`
			PredictedPoints   float64 `json:"PredictedPoints"`
		} `json:"18"`
		Field11 struct {
			RacingNumber      string  `json:"RacingNumber"`
			CurrentPosition   int     `json:"CurrentPosition"`
			PredictedPosition int     `json:"PredictedPosition"`
			CurrentPoints     float64 `json:"CurrentPoints"`
			PredictedPoints   float64 `json:"PredictedPoints"`
		} `json:"7"`
		Field12 struct {
			RacingNumber      string  `json:"RacingNumber"`
			CurrentPosition   int     `json:"CurrentPosition"`
			PredictedPosition int     `json:"PredictedPosition"`
			CurrentPoints     float64 `json:"CurrentPoints"`
			PredictedPoints   float64 `json:"PredictedPoints"`
		} `json:"99"`
		Field13 struct {
			RacingNumber      string  `json:"RacingNumber"`
			CurrentPosition   int     `json:"CurrentPosition"`
			PredictedPosition int     `json:"PredictedPosition"`
			CurrentPoints     float64 `json:"CurrentPoints"`
			PredictedPoints   float64 `json:"PredictedPoints"`
		} `json:"31"`
		Field14 struct {
			RacingNumber      string  `json:"RacingNumber"`
			CurrentPosition   int     `json:"CurrentPosition"`
			PredictedPosition int     `json:"PredictedPosition"`
			CurrentPoints     float64 `json:"CurrentPoints"`
			PredictedPoints   float64 `json:"PredictedPoints"`
		} `json:"63"`
		Field15 struct {
			RacingNumber      string  `json:"RacingNumber"`
			CurrentPosition   int     `json:"CurrentPosition"`
			PredictedPosition int     `json:"PredictedPosition"`
			CurrentPoints     float64 `json:"CurrentPoints"`
			PredictedPoints   float64 `json:"PredictedPoints"`
		} `json:"5"`
		Field16 struct {
			RacingNumber      string  `json:"RacingNumber"`
			CurrentPosition   int     `json:"CurrentPosition"`
			PredictedPosition int     `json:"PredictedPosition"`
			CurrentPoints     float64 `json:"CurrentPoints"`
			PredictedPoints   float64 `json:"PredictedPoints"`
		} `json:"47"`
		Field17 struct {
			RacingNumber      string  `json:"RacingNumber"`
			CurrentPosition   int     `json:"CurrentPosition"`
			PredictedPosition int     `json:"PredictedPosition"`
			CurrentPoints     float64 `json:"CurrentPoints"`
			PredictedPoints   float64 `json:"PredictedPoints"`
		} `json:"10"`
		Field18 struct {
			RacingNumber      string  `json:"RacingNumber"`
			CurrentPosition   int     `json:"CurrentPosition"`
			PredictedPosition int     `json:"PredictedPosition"`
			CurrentPoints     float64 `json:"CurrentPoints"`
			PredictedPoints   float64 `json:"PredictedPoints"`
		} `json:"6"`
		Field19 struct {
			RacingNumber      string  `json:"RacingNumber"`
			PredictedPosition int     `json:"PredictedPosition"`
			PredictedPoints   float64 `json:"PredictedPoints"`
		} `json:"14"`
		Field20 struct {
			RacingNumber      string  `json:"RacingNumber"`
			PredictedPosition int     `json:"PredictedPosition"`
			PredictedPoints   float64 `json:"PredictedPoints"`
		} `json:"9"`
	} `json:"Drivers"`
	Teams struct {
		Mercedes struct {
			TeamName          string  `json:"TeamName"`
			CurrentPosition   int     `json:"CurrentPosition"`
			PredictedPosition int     `json:"PredictedPosition"`
			CurrentPoints     float64 `json:"CurrentPoints"`
			PredictedPoints   float64 `json:"PredictedPoints"`
		} `json:"Mercedes"`
		RedBullRacingHonda struct {
			TeamName          string  `json:"TeamName"`
			CurrentPosition   int     `json:"CurrentPosition"`
			PredictedPosition int     `json:"PredictedPosition"`
			CurrentPoints     float64 `json:"CurrentPoints"`
			PredictedPoints   float64 `json:"PredictedPoints"`
		} `json:"Red Bull Racing Honda"`
		McLarenMercedes struct {
			TeamName          string  `json:"TeamName"`
			CurrentPosition   int     `json:"CurrentPosition"`
			PredictedPosition int     `json:"PredictedPosition"`
			CurrentPoints     float64 `json:"CurrentPoints"`
			PredictedPoints   float64 `json:"PredictedPoints"`
		} `json:"McLaren Mercedes"`
		Ferrari struct {
			TeamName          string  `json:"TeamName"`
			CurrentPosition   int     `json:"CurrentPosition"`
			PredictedPosition int     `json:"PredictedPosition"`
			CurrentPoints     float64 `json:"CurrentPoints"`
			PredictedPoints   float64 `json:"PredictedPoints"`
		} `json:"Ferrari"`
		AlphaTauriHonda struct {
			TeamName          string  `json:"TeamName"`
			CurrentPosition   int     `json:"CurrentPosition"`
			PredictedPosition int     `json:"PredictedPosition"`
			CurrentPoints     float64 `json:"CurrentPoints"`
			PredictedPoints   float64 `json:"PredictedPoints"`
		} `json:"AlphaTauri Honda"`
		AstonMartinMercedes struct {
			TeamName          string  `json:"TeamName"`
			CurrentPosition   int     `json:"CurrentPosition"`
			PredictedPosition int     `json:"PredictedPosition"`
			CurrentPoints     float64 `json:"CurrentPoints"`
			PredictedPoints   float64 `json:"PredictedPoints"`
		} `json:"Aston Martin Mercedes"`
		AlfaRomeoRacingFerrari struct {
			TeamName          string  `json:"TeamName"`
			CurrentPosition   int     `json:"CurrentPosition"`
			PredictedPosition int     `json:"PredictedPosition"`
			CurrentPoints     float64 `json:"CurrentPoints"`
			PredictedPoints   float64 `json:"PredictedPoints"`
		} `json:"Alfa Romeo Racing Ferrari"`
		AlpineRenault struct {
			TeamName          string  `json:"TeamName"`
			CurrentPosition   int     `json:"CurrentPosition"`
			PredictedPosition int     `json:"PredictedPosition"`
			CurrentPoints     float64 `json:"CurrentPoints"`
			PredictedPoints   float64 `json:"PredictedPoints"`
		} `json:"Alpine Renault"`
		WilliamsMercedes struct {
			TeamName          string  `json:"TeamName"`
			CurrentPosition   int     `json:"CurrentPosition"`
			PredictedPosition int     `json:"PredictedPosition"`
			CurrentPoints     float64 `json:"CurrentPoints"`
			PredictedPoints   float64 `json:"PredictedPoints"`
		} `json:"Williams Mercedes"`
		HaasFerrari struct {
			TeamName          string  `json:"TeamName"`
			CurrentPosition   int     `json:"CurrentPosition"`
			PredictedPosition int     `json:"PredictedPosition"`
			CurrentPoints     float64 `json:"CurrentPoints"`
			PredictedPoints   float64 `json:"PredictedPoints"`
		} `json:"Haas Ferrari"`
	} `json:"Teams"`
}
