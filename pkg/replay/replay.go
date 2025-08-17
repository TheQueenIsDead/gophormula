package replay

import (
	"bufio"
	"bytes"
	"errors"
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
		//defer file.Close()
	}

	return nil
}

func (r *Replayer) broadcast(message any) {
	for _, subscriber := range r.subscribers {
		ch := *subscriber
		ch <- message
	}
}

func (r *Replayer) Start() {
	for _, file := range r.files {
		scanner := bufio.NewScanner(&file)
		for scanner.Scan() {
			line := scanner.Text()
			r.broadcast(line)
			// FIXME: Need to remove this in place of logic that pauses at appropriate times to simulate replay
			time.Sleep(1 * time.Millisecond)
		}
	}
}

func (r *Replayer) Subscribe() <-chan any {
	ch := make(chan any, 64)
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

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
