package messages

import (
	"encoding/json"
)

type Feed struct {
	KeyFramePath string `json:"KeyFramePath"`
	StreamPath   string `json:"StreamPath"`
}
type Index struct {
	Feeds struct {
		SessionInfo            Feed `json:"SessionInfo"`
		ArchiveStatus          Feed `json:"ArchiveStatus"`
		TrackStatus            Feed `json:"TrackStatus"`
		SessionData            Feed `json:"SessionData"`
		ContentStreams         Feed `json:"ContentStreams"`
		ChampionshipPrediction Feed `json:"ChampionshipPrediction"`
		AudioStreams           Feed `json:"AudioStreams"`
		TimingDataF1           Feed `json:"TimingDataF1"`
		TimingData             Feed `json:"TimingData"`
		DriverList             Feed `json:"DriverList"`
		SPFeed                 Feed `json:"SPFeed"`
		LapSeries              Feed `json:"LapSeries"`
		TopThree               Feed `json:"TopThree"`
		TimingAppData          Feed `json:"TimingAppData"`
		TimingStats            Feed `json:"TimingStats"`
		CarDataZ               Feed `json:"CarData.z"`
		PositionZ              Feed `json:"Position.z"`
		ExtrapolatedClock      Feed `json:"ExtrapolatedClock"`
		TyreStintSeries        Feed `json:"TyreStintSeries"`
		DriverRaceInfo         Feed `json:"DriverRaceInfo"`
		LapCount               Feed `json:"LapCount"`
		SessionStatus          Feed `json:"SessionStatus"`
		Heartbeat              Feed `json:"Heartbeat"`
		WeatherData            Feed `json:"WeatherData"`
		WeatherDataSeries      Feed `json:"WeatherDataSeries"`
		TeamRadio              Feed `json:"TeamRadio"`
		TlaRcm                 Feed `json:"TlaRcm"`
		RaceControlMessages    Feed `json:"RaceControlMessages"`
		CurrentTyres           Feed `json:"CurrentTyres"`
		DriverScore            Feed `json:"DriverScore"`
		PitLaneTimeCollection  Feed `json:"PitLaneTimeCollection"`
	} `json:"Feeds"`
}

func (i *Index) String() string {
	buf, _ := json.MarshalIndent(i, "", "  ")
	return string(buf)
}

func (i *Index) GetFeeds() []Feed {
	return []Feed{
		i.Feeds.SessionInfo,
		i.Feeds.ArchiveStatus,
		i.Feeds.TrackStatus,
		i.Feeds.SessionData,
		i.Feeds.ContentStreams,
		i.Feeds.ChampionshipPrediction,
		i.Feeds.AudioStreams,
		i.Feeds.TimingDataF1,
		i.Feeds.TimingData,
		i.Feeds.DriverList,
		i.Feeds.SPFeed,
		i.Feeds.LapSeries,
		i.Feeds.TopThree,
		i.Feeds.TimingAppData,
		i.Feeds.TimingStats,
		i.Feeds.CarDataZ,
		i.Feeds.PositionZ,
		i.Feeds.ExtrapolatedClock,
		i.Feeds.TyreStintSeries,
		i.Feeds.DriverRaceInfo,
		i.Feeds.LapCount,
		i.Feeds.SessionStatus,
		i.Feeds.Heartbeat,
		i.Feeds.WeatherData,
		i.Feeds.WeatherDataSeries,
		i.Feeds.TeamRadio,
		i.Feeds.TlaRcm,
		i.Feeds.RaceControlMessages,
		i.Feeds.CurrentTyres,
		i.Feeds.DriverScore,
		i.Feeds.PitLaneTimeCollection,
	}
}
