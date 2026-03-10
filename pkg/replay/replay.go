package replay

import (
	"bufio"
	"errors"
	"fmt"
	"gophormula/pkg/livetiming"
	"io"
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

// peekFirstTimestamp reads the first line with a valid timestamp from f without
// consuming the file — the caller must seek f back to 0 after calling this.
func peekFirstTimestamp(f *os.File) *time.Time {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		ts, _, _ := livetiming.ExtractReplayData(scanner.Text())
		if ts != nil {
			return ts
		}
	}
	return nil
}

func (r *Replayer) Start() error {
	// Pre-scan every file to find the session origin — the earliest timestamp
	// across all files. All files are rewound afterwards so playback starts
	// from the beginning. This synchronises files that begin at different
	// points within the session onto a single real-time axis.
	var sessionOrigin *time.Time
	for i := range r.files {
		ts := peekFirstTimestamp(&r.files[i])
		if _, err := r.files[i].Seek(0, io.SeekStart); err != nil {
			return fmt.Errorf("rewinding %s: %w", r.files[i].Name(), err)
		}
		if ts == nil {
			continue
		}
		if sessionOrigin == nil || ts.Before(*sessionOrigin) {
			t := *ts
			sessionOrigin = &t
		}
	}

	wallStart := time.Now()

	for i := range r.files {
		go func(f *os.File) {
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				ts, msg, err := livetiming.ExtractReplayData(line)
				if err != nil {
					continue
				}
				// Sleep until the wall-clock time that corresponds to this
				// message's session timestamp, keeping all files in sync.
				if ts != nil && sessionOrigin != nil {
					offset := ts.Sub(*sessionOrigin)
					if d := time.Until(wallStart.Add(offset)); d > 0 {
						time.Sleep(d)
					}
				}
				_ = livetiming.ParseJSON(msg)
				r.broadcast(line)
			}
		}(&r.files[i])
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
