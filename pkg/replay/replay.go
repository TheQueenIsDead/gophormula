package replay

import (
	"bufio"
	"errors"
	"gophormula/pkg/livetiming"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Replayer struct {
	files       []os.File
	subscribers []*chan any
}

func New() *Replayer {
	return &Replayer{}
}

func (r *Replayer) ParseGlob(glob string) error {

	matches, err := filepath.Glob(glob)
	if err != nil {
		log.Println("error matching glob: ", err)
		return err
	}
	if len(matches) == 0 {
		log.Println("no matching files")
		return err
	}

	for _, match := range matches {
		matchIsHiddenFile := strings.HasSuffix(match, ".")
		matchIsIndex := match == "Index.json"
		if matchIsHiddenFile || matchIsIndex {
			continue
		}

		log.Println("processing", filepath.Base(match))

		file, err := os.Open(match)
		if err != nil {
			log.Fatal(err)
		}
		r.files = append(r.files, *file)
	}

	return nil
}

func (r *Replayer) broadcast(message any) {
	for _, subscriber := range r.subscribers {
		ch := *subscriber
		ch <- message
	}
}

// FIXME: Need to account for stream data starting at different times in different files
func (r *Replayer) Start() error {
	for _, file := range r.files {
		//topic := filepath.Base(file.Name())
		go func() {
			scanner := bufio.NewScanner(&file)
			for scanner.Scan() {
				line := scanner.Text()
				// TODO: FIXME
				_, msg, err := livetiming.ExtractReplayData(line)
				if err != nil {
					return
				}
				_ = livetiming.ParseJSON(msg)
				r.broadcast(line)
				// FIXME: Need to remove this in place of logic that pauses at appropriate times to simulate replay
				time.Sleep(1 * time.Second)
			}
		}()
	}
	return nil
}

func (r *Replayer) Subscribe() <-chan any {
	ch := make(chan any)
	r.subscribers = append(r.subscribers, &ch)
	return ch
}

func (r *Replayer) StartAndSubscribe() <-chan any {
	go r.Start()
	return r.Subscribe()
}

func (r *Replayer) Close() error {
	var errs []error
	for _, file := range r.files {
		err := file.Close()
		if err != nil {
			errs = append(errs, err)
			log.Printf("error closing file: %x\n", file.Name())
		}
	}
	return errors.Join(errs...)
}
