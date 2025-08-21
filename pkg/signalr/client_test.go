package signalr

import (
	"testing"
)

// Hodge-podge way of testing...But there's an old echo server container
// that can be spun up with the following cheeky beaky:
//
//	docker run --rm -p 6969:5000 ghcr.io/equinor/signalr-echo-server
func TestClient(t *testing.T) {

	opts := []Option{
		WithURL("http://localhost:6969/echo"),
		WithAck(false),
		WithVersion(1),
	}
	client := NewClient(opts...)
	err := client.Connect()
	if err != nil {
		t.Error(err)
	}
}
