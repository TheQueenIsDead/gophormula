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
	for {
		select {
		case m := <-ch:
			log.Println("Received:", m)
			switch m.(type) {
			case signalr.Message:
				message := m.(signalr.Message)
				log.Println(message)
			case []byte:
				buf := m.([]byte)
				_, err := livetiming.Classify(buf)
				if err != nil {
					log.Fatal(err)
				}
			case string:
				s := m.(string)
				_, err := livetiming.Classify([]byte(s))
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}

}
