package livetiming

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"reflect"
	"strings"
	"time"
)

var messageRegistry = make(map[string]reflect.Type)

func init() {
	// Register all known message types
	messageRegistry["Heartbeat"] = reflect.TypeOf(Heartbeat{})
	messageRegistry["CarData"] = reflect.TypeOf(CarData{})
	messageRegistry["Position"] = reflect.TypeOf(PositionData{})
	messageRegistry["SessionInfo"] = reflect.TypeOf(SessionInfo{})
	messageRegistry["TimingData"] = reflect.TypeOf(TimingData{})
	messageRegistry["TopThree"] = reflect.TypeOf(TopThree{})
	messageRegistry["TimingStats"] = reflect.TypeOf(TimingStats{})
	messageRegistry["TimingAppData"] = reflect.TypeOf(TimingAppData{})
	messageRegistry["WeatherData"] = reflect.TypeOf(WeatherData{})
	messageRegistry["TrackStatus"] = reflect.TypeOf(TrackStatus{})
	messageRegistry["DriverList"] = reflect.TypeOf(DriverList{})
	messageRegistry["RaceControlMessages"] = reflect.TypeOf(RaceControlMessages{})
	messageRegistry["SessionData"] = reflect.TypeOf(SessionData{})
	messageRegistry["LapCount"] = reflect.TypeOf(LapCount{})
}

// Decompress takes a byte slice, decodes it from base64, and decompresses it using flate.
func Decompress(data []byte) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, err
	}
	r := flate.NewReader(bytes.NewReader(decoded))
	defer func(r io.ReadCloser) {
		err := r.Close()
		if err != nil {
			log.Printf("failed to close reader: %v\n", err)
		}
	}(r)
	var out bytes.Buffer
	if _, err := out.ReadFrom(r); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func ParseJSON(msg json.RawMessage) []any {
	if len(msg) == 0 {
		return nil
	}

	// M format: array of SignalR hub invocations sent as incremental updates.
	// e.g. [{"H":"Streaming","M":"feed","A":["TopicName", data, "timestamp"]}]
	if msg[0] == '[' {
		var invocations []struct {
			A []json.RawMessage `json:"A"`
		}
		if err := json.Unmarshal(msg, &invocations); err != nil {
			log.Printf("error unmarshalling M invocations: %v", err)
			return nil
		}
		var results []any
		for _, inv := range invocations {
			if len(inv.A) < 2 {
				log.Printf("invalid invocation: expected at least 2 arguments, got %d", len(inv.A))
				continue
			}
			var topic string
			if err := json.Unmarshal(inv.A[0], &topic); err != nil {
				log.Printf("error unmarshalling topic from invocation: %v", err)
				continue
			}
			parsed, err := Parse(topic, inv.A[1])
			if err != nil {
				log.Printf("error parsing topic %s: %v", topic, err)
				continue
			}
			results = append(results, parsed)
		}
		return results
	}

	// R format: full state snapshot, a map of topic -> data.
	// e.g. {"Heartbeat": {...}, "CarData.z": "base64..."}
	var snapshot map[string]json.RawMessage
	if err := json.Unmarshal(msg, &snapshot); err != nil {
		log.Printf("error unmarshalling R snapshot: %v", err)
		return nil
	}
	var results []any
	for topic, data := range snapshot {
		parsed, err := Parse(topic, data)
		if err != nil {
			log.Printf("error parsing topic %s: %v", topic, err)
			continue
		}
		results = append(results, parsed)
	}
	return results
}

func Parse(topic string, data []byte) (any, error) {
	// If topic ends with .z, decompress
	if strings.HasSuffix(topic, ".z") {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return nil, fmt.Errorf("error unquoting compressed data: %w", err)
		}
		var err error
		data, err = Decompress([]byte(s))
		if err != nil {
			return nil, fmt.Errorf("error decompressing data: %w", err)
		}
		topic = strings.TrimSuffix(topic, ".z")
	}

	if t, ok := messageRegistry[topic]; ok {
		// Create a new instance of the message type
		v := reflect.New(t).Interface()
		if err := json.Unmarshal(data, v); err != nil {
			return nil, fmt.Errorf("error unmarshalling %s: %w", topic, err)
		}
		return v, nil
	}

	return nil, fmt.Errorf("unknown message topic: %s", topic)
}

func ExtractReplayData(line string) (*time.Time, json.RawMessage, error) {

	// Remove UTF8 BOM if present at the start of the line
	trimmed := bytes.TrimPrefix([]byte(line), []byte{0xEF, 0xBB, 0xBF})
	line = string(trimmed)

	// Check for a timestamp and remove if so
	var timestamp *time.Time
	if len(line) > 12 {
		ts := line[:12]
		t, err := time.Parse("15:04:05.999", ts)
		if err == nil {
			line = line[12:]
			timestamp = &t
		}
	}

	line = strings.TrimSpace(line)

	// Compressed data
	if !strings.HasPrefix(line, "{") {
		decompressed, err := Decompress([]byte(line))
		if err != nil {
			return timestamp, nil, err
		}
		return timestamp, decompressed, nil
	}

	// Non compressed data - post sanitization
	return timestamp, json.RawMessage(line), nil
}
