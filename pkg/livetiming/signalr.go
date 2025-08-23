package livetiming

type Topic string

func (t Topic) String() string {
	return string(t)
}

var (
	Heartbeat           Topic = "Heartbeat"
	Cardata             Topic = "Cardata.z"
	Position            Topic = "Position.z"
	ExtrapolatedClock   Topic = "ExtrapolatedClock"
	TopThree            Topic = "TopThree"
	RcmSeries           Topic = "RcmSeries"
	TimingStats         Topic = "TimingStats"
	TimingAppData       Topic = "TimingAppData"
	WeatherData         Topic = "WeatherData"
	TrackStatus         Topic = "TrackStatus"
	DriverList          Topic = "DriverList"
	RaceControlMessages Topic = "RaceControlMessages"
	SessionInfo         Topic = "SessionInfo"
	SessionData         Topic = "SessionData"
	LapCount            Topic = "LapCount"
	TimingData          Topic = "TimingData"
)

var (
	Topics = []Topic{
		Heartbeat, Cardata, Position,
		ExtrapolatedClock, TopThree, RcmSeries,
		TimingStats, TimingAppData, WeatherData,
		TrackStatus, DriverList, RaceControlMessages,
		SessionInfo, SessionData, LapCount,
		TimingData,
	}
)

func AllTopics() []string {
	var topics []string
	for _, t := range Topics {
		topics = append(topics, t.String())
	}
	return topics
}
