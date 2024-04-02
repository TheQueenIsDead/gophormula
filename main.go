package main

import "golang.org/x/net/websocket"

/*
self.topics = ["Heartbeat", "CarData.z", "Position.z",
"ExtrapolatedClock", "TopThree", "RcmSeries",
"TimingStats", "TimingAppData",
"WeatherData", "TrackStatus", "DriverList",
"RaceControlMessages", "SessionInfo",
"SessionData", "LapCount", "TimingData"]
*/
func main() {

	type T struct {
		Msg   string
		Count int
	}

	// receive JSON type T
	var data T
	websocket.JSON.Receive(ws, &data)
	websocket

	// send JSON type T
	websocket.JSON.Send(ws, data)
	//
	//req := http.Request{
	//	Method:     "",
	//	URL:        &url.URL{Scheme: "http", Host: "livetiming.formula1.com", Path: "/signalr"},
	//	Proto:      "",
	//	ProtoMajor: 0,
	//	ProtoMinor: 0,
	//	Header: map[string][]string{
	//		"User-agent":      {"BestHTTP"},
	//		"Accept-Encoding": {"gzip, identity"},
	//		"Connection":      {"keep-alive, Upgrade"},
	//	},
	//	Body:             nil,
	//	GetBody:          nil,
	//	ContentLength:    0,
	//	TransferEncoding: nil,
	//	Close:            false,
	//	Host:             "",
	//	Form:             nil,
	//	PostForm:         nil,
	//	MultipartForm:    nil,
	//	Trailer:          nil,
	//	RemoteAddr:       "",
	//	RequestURI:       "",
	//	TLS:              nil,
	//	Cancel:           nil,
	//	Response:         nil,
	//}
	//
	//client := http.Client{
	//	Transport:     nil,
	//	CheckRedirect: nil,
	//	Jar:           nil,
	//	Timeout:       0,
	//}
	//
	//res, err := client.Do(&req)
	//fmt.Print(err)
	//fmt.Print(res)
}
