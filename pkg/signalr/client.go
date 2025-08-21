package signalr

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
)

const (
	NegotiationPath = "/negotiate"

	// TODO: This may need to be used to track separate messages on the signalr transport:
	// https://github.com/dotnet/aspnetcore/blob/main/src/SignalR/docs/specs/HubProtocol.md#overview
	RecordSeparator = "\x1E"
)

type Client struct {
	url     string
	base    string
	version uint
	ack     bool
	token   string
}

type Option func(*Client)

func WithURL(url string) Option {
	return func(c *Client) {
		c.url = url
	}
}

func WithVersion(version uint) Option {
	return func(c *Client) {
		c.version = version
	}
}

func WithAck(ack bool) Option {
	return func(c *Client) {
		c.ack = ack
	}
}

func NewClient(opts ...Option) *Client {
	client := &Client{
		url:     "http://localhost:8080/signalr",
		version: 1,
		ack:     false,
	}

	for _, opt := range opts {
		opt(client)
	}
	return client
}

func (c *Client) Connect() error {

	uri, err := url.Parse(c.url)
	if err != nil {
		return err
	}

	negotiationResponse, err := c.negotiate(*(uri.JoinPath(NegotiationPath)))
	if err != nil {
		return err
	}
	fmt.Println(negotiationResponse)

	//if len(negotiationResponse.AvailableTransports) == 0 {
	//	return errors.New("no available transports")
	//}

	//var transport AvailableTransport
	//var transportPriority = math.MinInt32
	//for _, available := range negotiationResponse.AvailableTransports {
	//	priority := -1
	//	switch available.Transport {
	//	case "LongPolling":
	//		priority = 0
	//	case "ServerSentEvents":
	//		priority = 1
	//	case "WebSockets":
	//		priority = 2
	//	}
	//	if priority > transportPriority {
	//		transport = available
	//		transportPriority = priority
	//	}
	//}

	//fmt.Println("Chosen transport:", transport)
	fmt.Println("Chosen transport: websocket")

	// Initialise websocket
	u := url.URL{
		Scheme: "wss", // try wss instead of ws
		Host:   "livetiming.formula1.com",
		Path:   "/signalr/connect",
	}
	uri.Scheme = "wss"

	q := u.Query()
	q.Set("transport", "webSockets")
	q.Set("connectionToken", negotiationResponse.ConnectionToken)
	u.RawQuery = q.Encode()

	conn, resp, err := websocket.DefaultDialer.Dial(
		u.String(),
		//http.Header{"User-Agent": []string{"Go-SignalR-Client"}},
		http.Header{},
	)
	if err != nil {
		log.Fatalf("dial error: %v (status %v)", err, resp.Status)
	}
	defer conn.Close()
	log.Println("Connected!")

	// Send a HandshakeRequest
	b, _ := json.Marshal(HandshakeRequest{
		Protocol: "json",
		Version:  1,
	})
	// append record separator (0x1E)
	b = append(b, 0x1E)
	conn.WriteMessage(websocket.TextMessage, b)

	// Read server response
	_, msg, _ := conn.ReadMessage()
	// should be {} + 0x1E if successful
	fmt.Println(string(msg))

	fmt.Println("Reading messages...")
	buf := ""
	for {
		n, message, _ := conn.ReadMessage()
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

	return nil
}

func (c *Client) negotiate(url url.URL) (*NegotiationResponse, error) {

	url.Query().Add("version", fmt.Sprint(c.version))
	if c.ack {
		url.Query().Add("useAck", fmt.Sprint(c.ack))
	}

	res, err := http.Post(url.String(), "application/json", nil)
	if err != nil || res.StatusCode != 200 {
		return nil, err
	}

	var negotiationResponse NegotiationResponse
	err = json.NewDecoder(res.Body).Decode(&negotiationResponse)
	if err != nil {
		return nil, err
	}

	return &negotiationResponse, nil
}
