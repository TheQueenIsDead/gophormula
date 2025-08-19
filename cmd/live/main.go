package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
)

const (
	PROTOCOL       = "https://"
	FQDN           = "livetiming.formula1.com"
	BASE_PATH      = "/signalr"
	NEGOTIATE_PATH = "/negotiate"
)

type NegotiationResponse struct {
	Url                     string  `json:"Url"`
	ConnectionToken         string  `json:"ConnectionToken"`
	ConnectionId            string  `json:"ConnectionId"`
	KeepAliveTimeout        float64 `json:"KeepAliveTimeout"`
	DisconnectTimeout       float64 `json:"DisconnectTimeout"`
	ConnectionTimeout       float64 `json:"ConnectionTimeout"`
	TryWebSockets           bool    `json:"TryWebSockets"`
	ProtocolVersion         string  `json:"ProtocolVersion"`
	TransportConnectTimeout float64 `json:"TransportConnectTimeout"`
	LongPollDelay           float64 `json:"LongPollDelay"`
}

func (r *NegotiationResponse) String() string {
	buf, _ := json.MarshalIndent(r, "", "  ")
	return string(buf)
}

func main() {

	// Build negotiation
	negotiate, err := url.JoinPath(PROTOCOL, FQDN, BASE_PATH, NEGOTIATE_PATH)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(negotiate)

	// Start
	res, err := http.Get(negotiate)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// Marshal the request body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var negotiationResponse NegotiationResponse
	err = json.Unmarshal(body, &negotiationResponse)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(negotiationResponse.String())

	// Choose a protocol to negotiate
	// FIXME: Defaulted to websocket, but would be good to negotiate later

	// Initialise websocket
	u := url.URL{
		Scheme: "wss", // try wss instead of ws
		Host:   "livetiming.formula1.com",
		Path:   "/signalr/connect",
	}
	q := u.Query()
	q.Set("transport", "webSockets")
	q.Set("connectionToken", negotiationResponse.ConnectionToken)
	u.RawQuery = q.Encode()

	c, resp, err := websocket.DefaultDialer.Dial(
		u.String(),
		//http.Header{"User-Agent": []string{"Go-SignalR-Client"}},
		http.Header{},
	)
	if err != nil {
		log.Fatalf("dial error: %v (status %v)", err, resp.Status)
	}
	defer c.Close()
	log.Println("Connected!")

	// Send a HandshakeRequest
	// Handshake request
	handshake := map[string]interface{}{
		"protocol": "json",
		"version":  1,
	}

	b, _ := json.Marshal(handshake)
	// append record separator (0x1E)
	b = append(b, 0x1E)
	c.WriteMessage(websocket.TextMessage, b)

	// Read server response
	_, msg, _ := c.ReadMessage()
	// should be {} + 0x1E if successful
	fmt.Println(string(msg))

	fmt.Println("Reading messages...")
	buf := ""
	for {
		n, message, _ := c.ReadMessage()
		fmt.Printf("Read %d bytes: %s\n", n, message)
		buf += string(message)
		parts := strings.Split(buf, string(rune(0x1E)))
		for _, p := range parts[:len(parts)-1] {
			if p == "" {
				continue
			}
			var m map[string]interface{}
			json.Unmarshal([]byte(p), &m)
			fmt.Printf("Got message: %+v\n", m)
		}
		buf = parts[len(parts)-1] // leftover
	}

}
