package main

import (
	"gophormula/pkg/livetiming"
	"gophormula/pkg/signalr"
	"log"
)

func main() {
	client := signalr.NewClient(
		signalr.WithURL("https://livetiming.formula1.com/signalr"),
	)

	ch, err := client.Connect(
		[]signalr.Hub{"Streaming"},
		livetiming.AllTopics(),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Waiting for messages...")
	for {
		select {
		case msg := <-ch:
			data := livetiming.ParseJSON(msg.Data())
			switch v := data.(type) {
			case *livetiming.Heartbeat:
				log.Printf("[Heartbeat] Received at %v", v)
			case *livetiming.CarData:
				log.Printf("[CarData] Received for %d cars", len(v.Entries))
			case *livetiming.PositionData:
				log.Printf("[PositionData] Received")
			case *livetiming.SessionInfo:
				log.Printf("[SessionInfo] %s - %s", v.Meeting.Name, v.Name)
			case *livetiming.TimingData:
				log.Printf("[TimingData] Received for %d lines", len(v.Lines))
			case *livetiming.TopThree:
				log.Printf("[TopThree] Received with %d drivers", len(v.Lines))
			case *livetiming.TimingStats:
				log.Printf("[TimingStats] Received for %d lines", len(v.Lines))
			case *livetiming.TimingAppData:
				log.Printf("[TimingAppData] Received for %d lines", len(v.Lines))
			case *livetiming.WeatherData:
				log.Printf("[WeatherData] Air Temp: %s, Track Temp: %s", v.AirTemp, v.TrackTemp)
			case *livetiming.TrackStatus:
				log.Printf("[TrackStatus] Status: %s - %s", v.Status, v.Message)
			case *livetiming.DriverList:
				log.Printf("[DriverList] Received for %d drivers", len(*v))
			case *livetiming.RaceControlMessages:
				log.Printf("[RaceControl] Received %d new livetiming", len(v.Messages))
			case *livetiming.SessionData:
				log.Printf("[SessionData] Received with %d data points", len(v.Series))
			case *livetiming.LapCount:
				log.Printf("[LapCount] Lap %d/%d", v.CurrentLap, v.TotalLaps)
			default:
				log.Printf("Received unknown message type: %T: %s", data, v)
			}
		}
	}
}
