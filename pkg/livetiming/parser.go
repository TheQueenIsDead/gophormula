package livetiming

import (
	"bufio"
	"bytes"
	"compress/flate"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"maps"
	"slices"
	"strings"
	"time"
)

const (
	// IntervalTimeFormat is stdHour:stdZeroMinute:stdZeroSecond.stdFracSecond9 to 3 decimal places
	// Specifics about the date format can be found at https://go.dev/src/time/format.go
	IntervalTimeFormat = "15:04:05.999"
)

func Parse(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		_, l, err := ParseLine(scanner.Text())
		log.Printf("%v\n", l)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func ParseLine(line string) (*time.Time, *string, error) {

	// Sanitise the line to remove incorrectly appended UTF-8 BOM's
	line = strings.ToValidUTF8(line, "")
	line = strings.Replace(line, "\xEF\xBB\xBF", "", -1)

	// Check if there is a timestamp at the beginning of the line, else use now
	twelve := line[:min(len(line), 12)]
	timestamp, tsErr := time.Parse(IntervalTimeFormat, twelve)
	if tsErr != nil {
		timestamp = time.Now()
	}

	// If the ts was valid, remove it from the line
	if tsErr == nil {
		line = line[min(len(line), 12):]
	}

	// If the start of the line is not JSON compliant, assume compressed
	if line[0] != '{' {

		compressed := strings.ReplaceAll(line, "\"", "")
		decoded, err := base64.StdEncoding.DecodeString(compressed)
		if err != nil {
			log.Printf("failed to decode base64 string: %v\n", err)
			log.Printf("line: %v\n", line)
			return nil, nil, err
		}

		r := flate.NewReader(bytes.NewReader(decoded))
		defer func(r io.ReadCloser) {
			err := r.Close()
			if err != nil {
				log.Printf("failed to close reader: %v\n", err)
				panic(err)
			}
		}(r)
		var out bytes.Buffer
		if _, err := out.ReadFrom(r); err != nil {
			log.Println("flate decompress error", err)
			return nil, nil, err
		}
		line = out.String()
	}

	// Attempt to json decode the payload
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(line), &payload); err != nil {
		log.Println(err)
		return nil, nil, err
	}

	// TODO: Once we know the keys, marshal into the corresponding struct. Currently WIP
	keySeq := maps.Keys(payload)
	keyArr := slices.Collect(keySeq)
	if len(keyArr) == 0 {
		return nil, nil, errors.New("no json keys found")
	}
	keys := strings.Join(keyArr, "-")

	return &timestamp, &keys, nil

}
