package livetiming

type LiveTimingFile string

const (
	LiveTimingFileSessionData            LiveTimingFile = "SessionData.json"                  // track + session status + lap count
	LiveTimingFileSessionInfo            LiveTimingFile = "SessionInfo.jsonStream"            // # more rnd
	LiveTimingFileArchiveStatus          LiveTimingFile = "ArchiveStatus.json"                // rnd=1880327548
	LiveTimingFileHeartBeat              LiveTimingFile = "HeartBeat.jsonStream"              // Probably time synchronization?
	LiveTimingFileAudioStreams           LiveTimingFile = "AudioStreams.jsonStream"           // # Link to audio commentary
	LiveTimingFileDriverList             LiveTimingFile = "DriverList.jsonStream"             // # Driver info and line story
	LiveTimingFileExtrapolatedClock      LiveTimingFile = "ExtrapolatedClock.jsonStream"      // Boolean
	LiveTimingFileRaceControlMessages    LiveTimingFile = "RaceControlMessages.jsonStream"    // Flags etc
	LiveTimingFileSessionStatus          LiveTimingFile = "SessionStatus.jsonStream"          // Start and finish times
	LiveTimingFileTeamRadio              LiveTimingFile = "TeamRadio.jsonStream"              // Links to team radios
	LiveTimingFileTimingAppData          LiveTimingFile = "TimingAppData.jsonStream"          // Tyres and laps (juicy)
	LiveTimingFileTimingStats            LiveTimingFile = "TimingStats.jsonStream"            // 'Best times/speed' useless
	LiveTimingFileTrackStatus            LiveTimingFile = "TrackStatus.jsonStream"            // SC, VSC and Yellow
	LiveTimingFileWeatherData            LiveTimingFile = "WeatherData.jsonStream"            // Temp, wind and rain
	LiveTimingFilePosition               LiveTimingFile = "Position.z.jsonStream"             // Coordinates, not GPS? (.z)
	LiveTimingFileCarData                LiveTimingFile = "CarData.z.jsonStream"              // Telemetry channels (.z)
	LiveTimingFileContentStreams         LiveTimingFile = "ContentStreams.jsonStream"         // Lap by lap feeds
	LiveTimingFileTimingData             LiveTimingFile = "TimingData.jsonStream"             // Gap to car ahead
	LiveTimingFileLapCount               LiveTimingFile = "LapCount.jsonStream"               // Lap counter
	LiveTimingFileChampionshipPrediction LiveTimingFile = "ChampionshipPrediction.jsonStream" // Points
	LiveTimingFileIndex                  LiveTimingFile = "Index.json"
)

//const (
//	SessionData            = "SessionData.json"                  // track + session status + lap count
//	SessionInfo            = "SessionInfo.jsonStream"            // # more rnd
//	ArchiveStatus          = "ArchiveStatus.json"                // rnd=1880327548
//	HeartBeat              = "HeartBeat.jsonStream"              // Probably time synchronization?
//	AudioStreams           = "AudioStreams.jsonStream"           // # Link to audio commentary
//	DriverList             = "DriverList.jsonStream"             // # Driver info and line story
//	ExtrapolatedClock      = "ExtrapolatedClock.jsonStream"      // Boolean
//	RaceControlMessages    = "RaceControlMessages.jsonStream"    // Flags etc
//	SessionStatus          = "SessionStatus.jsonStream"          // Start and finish times
//	TeamRadio              = "TeamRadio.jsonStream"              // Links to team radios
//	TimingAppData          = "TimingAppData.jsonStream"          // Tyres and laps (juicy)
//	TimingStats            = "TimingStats.jsonStream"            // 'Best times/speed' useless
//	TrackStatus            = "TrackStatus.jsonStream"            // SC, VSC and Yellow
//	WeatherData            = "WeatherData.jsonStream"            // Temp, wind and rain
//	Position               = "Position.z.jsonStream"             // Coordinates, not GPS? (.z)
//	CarData                = "CarData.z.jsonStream"              // Telemetry channels (.z)
//	ContentStreams         = "ContentStreams.jsonStream"         // Lap by lap feeds
//	TimingData             = "TimingData.jsonStream"             // Gap to car ahead
//	LapCount               = "LapCount.jsonStream"               // Lap counter
//	ChampionshipPrediction = "ChampionshipPrediction.jsonStream" // Points
//	Index                  = "Index.json"
//)
