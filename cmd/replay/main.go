package main

import (
	"fmt"
	"github.com/philippseith/signalr"
	"gophormula/pkg/replay"
	"log"
	"os"
	"path/filepath"
)

type ReplayHub struct {
}

func (r ReplayHub) Initialize(hubContext signalr.HubContext) {
	//TODO implement me
	panic("implement me")
}

func (r ReplayHub) OnConnected(connectionID string) {
	//TODO implement me
	panic("implement me")
}

func (r ReplayHub) OnDisconnected(connectionID string) {
	//TODO implement me
	panic("implement me")
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: replay <dir>")
	}

	dir := os.Args[1]

	replayer := replay.New()
	err := replayer.ParseGlob(filepath.Join(dir, "*"))
	if err != nil {
		log.Fatal(err)
	}
	ch := replayer.StartAndSubscribe()
	defer replayer.Close()

	for {
		select {
		case message := <-ch:
			fmt.Println(message)
		}
	}
}
