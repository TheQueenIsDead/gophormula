package main

import (
	"fmt"
	"gophormula/pkg/replay"
	"log"
	"os"
	"path/filepath"
)

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
