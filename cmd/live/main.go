package main

import (
	"fmt"
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
			fmt.Println("Received:", m)
		}
	}

}
