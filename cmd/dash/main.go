package main

import (
	"gophormula/pkg/dash"
	"log"
	"net/http"
	"os"
)

func main() {
	dataDir := "data"
	if len(os.Args) >= 2 {
		dataDir = os.Args[1]
	}

	hub := dash.NewHub(dataDir)

	mux := http.NewServeMux()
	mux.HandleFunc("/", hub.Index)
	mux.HandleFunc("/events", hub.Events)
	mux.HandleFunc("/replay", hub.ReplayHandler())

	log.Println("listening on :1234")
	if err := http.ListenAndServe(":1234", mux); err != nil {
		panic(err)
	}
}
