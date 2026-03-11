package livetiming

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// CircuitMap holds the track outline from the Multiviewer API.
// X and Y are parallel arrays of track coordinates (same space as Position.z data).
// Rotation is the suggested display rotation in degrees.
type CircuitMap struct {
	X        []float64 `json:"x"`
	Y        []float64 `json:"y"`
	Rotation float64   `json:"rotation"`
}

// FetchCircuitMap retrieves the circuit map from the Multiviewer API using the
// circuit key (from SessionInfo.Meeting.Circuit.Key) and the session year.
func FetchCircuitMap(circuitKey, year int) (*CircuitMap, error) {
	url := fmt.Sprintf("https://api.multiviewer.app/api/v1/circuits/%d/%d", circuitKey, year)
	resp, err := http.Get(url) //nolint:noctx
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("multiviewer API returned %d", resp.StatusCode)
	}
	var cm CircuitMap
	if err := json.NewDecoder(resp.Body).Decode(&cm); err != nil {
		return nil, err
	}
	return &cm, nil
}
