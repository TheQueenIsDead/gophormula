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

func ParseJSON(msg json.RawMessage) any {

	// Attempt to parse out JSON from the 'R' incremental update mode.
	var r map[string]json.RawMessage
	if err := json.Unmarshal(msg, &r); err == nil {
		for topic, data := range r {
			// TODO: Something!
			fmt.Println(topic, data)
		}
		return r
	}

	// Conversely, attempt to parse JSON from the standard 'M' message array.
	var m [][]json.RawMessage
	if err := json.Unmarshal(msg, &m); err == nil {
		for _, message := range m {
			if len(message) != 2 {
				log.Printf("invalid message format in M block: expected 2 elements, got %d", len(message))
				continue
			}
			var topic string
			if err := json.Unmarshal(message[0], &topic); err != nil {
				log.Printf("error unmarshalling topic from M block: %v", err)
				continue
			}
			// TODO: Something!
			fmt.Println(topic, message[1])
		}
		return m
	}
	return nil
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
