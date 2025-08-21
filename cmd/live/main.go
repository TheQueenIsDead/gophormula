package main

import (
	"gophormula/pkg/signalr"
	"log"
)

const (
	PROTOCOL       = "https://"
	FQDN           = "livetiming.formula1.com"
	BASE_PATH      = "/signalr"
	NEGOTIATE_PATH = "/negotiate"
)

func main() {

	client := signalr.NewClient(
		signalr.WithURL("https://livetiming.formula1.com/signalr"),
	)

	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}

}
