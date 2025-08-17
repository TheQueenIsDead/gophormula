package livetiming

import "time"

type WeatherDataPoint struct {
	Timestamp time.Time   `json:"Timestamp"`
	Weather   WeatherData `json:"Weather"`
}

type WeatherDataSeries struct {
	Series []WeatherDataPoint `json:"Series"`
}
