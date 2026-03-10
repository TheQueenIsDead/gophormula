package replay

import (
	"bufio"
	"bytes"
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

type fileEntry struct {
	file     os.File
	topic    string
	isStream bool // true for .jsonStream (line-by-line), false for .json (whole document)
}

// Message pairs a parsed livetiming value with the session timestamp from the
// stream file. Timestamp is nil for keyframe messages (no line timestamp).
type Message struct {
	Timestamp *time.Time
	Value     any
}

type Replayer struct {
	files       []fileEntry
	subscribers []*chan any
}

func New() *Replayer {
	return &Replayer{}
}

// topicFromFilename strips the .jsonStream or .json extension from a filename
// to derive the F1 topic name. e.g. "CarData.z.jsonStream" → "CarData.z".
// Returns the topic name and whether this is a stream file (vs a keyframe file).
func topicFromFilename(name string) (topic string, isStream bool) {
	base := filepath.Base(name)
	if strings.HasSuffix(base, ".jsonStream") {
		return strings.TrimSuffix(base, ".jsonStream"), true
	}
	return strings.TrimSuffix(base, ".json"), false
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
		base := filepath.Base(match)
		matchIsHiddenFile := strings.HasPrefix(base, ".")
		matchIsIndex := base == "Index.json"
		if matchIsHiddenFile || matchIsIndex {
			continue
		}

		log.Println("processing", filepath.Base(match))

		file, err := os.Open(match)
		if err != nil {
			log.Fatal(err)
		}
		topic, isStream := topicFromFilename(match)
		r.files = append(r.files, fileEntry{
			file:     *file,
			topic:    topic,
			isStream: isStream,
		})
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
		if !r.files[i].isStream {
			continue // keyframe files have no timestamps; skip in pre-scan
		}
		ts := peekFirstTimestamp(&r.files[i].file)
		if _, err := r.files[i].file.Seek(0, io.SeekStart); err != nil {
			return fmt.Errorf("rewinding %s: %w", r.files[i].file.Name(), err)
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
		if r.files[i].isStream {
			go func(f *os.File, topic string) {
				scanner := bufio.NewScanner(f)
				scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
				for scanner.Scan() {
					ts, msg, err := livetiming.ExtractReplayData(scanner.Text())
					if err != nil || msg == nil {
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
					parsed, err := livetiming.Parse(topic, msg)
					if err != nil {
						log.Printf("error parsing %s: %v", topic, err)
						continue
					}
					r.broadcast(Message{Timestamp: ts, Value: parsed})
				}
			}(&r.files[i].file, r.files[i].topic)
		} else {
			go func(f *os.File, topic string) {
				raw, err := io.ReadAll(f)
				if err != nil {
					log.Printf("error reading keyframe %s: %v", topic, err)
					return
				}
				// Strip UTF-8 BOM if present
				data := bytes.TrimPrefix(raw, []byte{0xEF, 0xBB, 0xBF})
				parsed, err := livetiming.Parse(topic, data)
				if err != nil {
					log.Printf("error parsing keyframe %s: %v", topic, err)
					return
				}
				r.broadcast(Message{Timestamp: nil, Value: parsed})
			}(&r.files[i].file, r.files[i].topic)
		}
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
	for _, entry := range r.files {
		err := entry.file.Close()
		if err != nil {
			errs = append(errs, err)
			log.Printf("error closing file: %x\n", entry.file.Name())
		}
	}
	return errors.Join(errs...)
}
