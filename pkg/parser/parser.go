package parser

import (
	"bufio"
	"bytes"
	"compress/flate"
	"encoding/base64"
	"fmt"
	"io"
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
		_, l, err := parseLine(scanner.Text())
		fmt.Println(*l)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

// parseLine takes a line from a *.z.jsonStream file and returns the decoded data.
// A line is expected to have the following structure:
// 00:00:03.712"base64encode(flateCompress(data))"
func parseLine(line string) (*time.Time, *string, error) {

	// Sanitise the line to remove incorrectly appended UTF-8 BOM's
	line = strings.ToValidUTF8(line, "")
	line = strings.Replace(line, "\xEF\xBB\xBF", "", -1)

	// Split the string on quotes to parse out the following parts
	// 0 - Interval
	// 1 - Message
	// 2 - Empty String
	parts := strings.Split(line, "\"")
	if len(parts) != 3 {
		return nil, nil, fmt.Errorf("invalid line: %s", line)
	}
	intervalData, messageData := parts[0], parts[1]

	interval, err := parseInterval(intervalData)
	if err != nil {
		fmt.Printf("%x\n", []byte(intervalData))
		fmt.Printf("%s\n", []byte(intervalData))
		fmt.Println(err)
		return nil, nil, fmt.Errorf("invalid interval: %s", intervalData)
	}

	message, err := parseMessage(messageData)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid message: %s", messageData)
	}

	return &interval, &message, nil
}

func parseInterval(interval string) (time.Time, error) {
	return time.Parse(IntervalTimeFormat, interval)
}

func parseMessage(message string) (string, error) {

	decoded, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 string: %v\n", err)
	}

	r := flate.NewReader(bytes.NewReader(decoded))
	defer func(r io.ReadCloser) {
		err := r.Close()
		if err != nil {
			panic(err)
		}
	}(r)
	var out bytes.Buffer
	if _, err := out.ReadFrom(r); err != nil {
		return "", fmt.Errorf("flate decompress error: %v\n", err)
	}

	return out.String(), nil
}
