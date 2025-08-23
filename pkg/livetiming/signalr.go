package livetiming

type Topic string

func (t Topic) String() string {
	return string(t)
}

var (
	TopicHeartbeat           Topic = "Heartbeat"
	TopicCardata             Topic = "Cardata.z"
	TopicPosition            Topic = "Position.z"
	TopicExtrapolatedClock   Topic = "ExtrapolatedClock"
	TopicTopThree            Topic = "TopThree"
	TopicRcmSeries           Topic = "RcmSeries"
	TopicTimingStats         Topic = "TimingStats"
	TopicTimingAppData       Topic = "TimingAppData"
	TopicWeatherData         Topic = "WeatherData"
	TopicTrackStatus         Topic = "TrackStatus"
	TopicDriverList          Topic = "DriverList"
	TopicRaceControlMessages Topic = "RaceControlMessages"
	TopicSessionInfo         Topic = "SessionInfo"
	TopicSessionData         Topic = "SessionData"
	TopicLapCount            Topic = "LapCount"
	TopicTimingData          Topic = "TimingData"
)

var (
	Topics = []Topic{
		TopicHeartbeat, TopicCardata, TopicPosition,
		TopicExtrapolatedClock, TopicTopThree, TopicRcmSeries,
		TopicTimingStats, TopicTimingAppData, TopicWeatherData,
		TopicTrackStatus, TopicDriverList, TopicRaceControlMessages,
		TopicSessionInfo, TopicSessionData, TopicLapCount,
		TopicTimingData,
	}
)

func AllTopics() []string {
	var topics []string
	for _, t := range Topics {
		topics = append(topics, t.String())
	}
	return topics
}
