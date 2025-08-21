package signalr

import (
	"encoding/json"
	"fmt"
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

func (c *Client) Connect() (chan interface{}, error) {

	uri, err := url.Parse(c.url)
	if err != nil {
		return nil, err
	}

	negotiationResponse, err := c.negotiate(*(uri.JoinPath(NegotiationPath)))
	if err != nil {
		return nil, err
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
	if !negotiationResponse.TryWebSockets {
		return nil, ErrWebsocketsUnsupported
	}
	fmt.Println("Chosen transport: websocket")

	transport := &WebsocketTransport{}
	conn, err := transport.Connect("livetiming.formula1.com/signalr", negotiationResponse.ConnectionToken, []Hub{{"Streaming"}})
	if err != nil {
		return nil, err
	}
	// defer conn.Close()

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
	fmt.Println("Init message....")
	fmt.Println(string(msg))

	// TODO: Send start request
	// Should receive: {Response: Started}
	/*
		» start – informs the server that transport started successfully
		Required parameters: transport, clientProtocol, connectionToken, connectionData (when using hubs)
		Optional parameters: queryString
		Sample request:

		http://host/signalr/start?transport=webSockets&clientProtocol=1.5&connectionToken=LkNk&connectionData=%5B%7B%22name%22%3A%22chat%22%7D%5D
		Sample response:
	*/

	// TODO: Invoke Subscription: hub.server.invoke("Subscribe", self.topics)
	//{"H":"chathub","M":"Send","A":["JS Client","Test message"],"I":0, "S":{"customProperty" : "abc"}}
	// Send a HandshakeRequest
	b, _ = json.Marshal(map[string]interface{}{
		"H": "Streaming",
		"M": "Subscribe",
		// It expects a single argument, an array of strings, but signalr uses an array of A, so nest the struct with [][]string
		"A": [][]string{{"Heartbeat", "CarData.z", "Position.z",
			"ExtrapolatedClock", "TopThree", "RcmSeries",
			"TimingStats", "TimingAppData",
			"WeatherData", "TrackStatus", "DriverList",
			"RaceControlMessages", "SessionInfo",
			"SessionData", "LapCount", "TimingData"}},
		"I": 1,
	})
	// append record separator (0x1E)
	b = append(b, 0x1E)
	conn.WriteMessage(websocket.TextMessage, b)
	fmt.Println("Sending", string(b))

	return c.read(conn), nil
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

func (c *Client) read(conn *websocket.Conn) chan interface{} {

	ch := make(chan interface{})

	go func() {
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
				// TODO: Parse keepalive ({}) which is sent every 10 seconds
			}
			buf = parts[len(parts)-1] // leftover

			ch <- string(message)
		}
	}()

	return ch
}

// TODO: Reconnect
/*
reconnect – sent to the server when the connection is lost and the client is reconnecting
Required parameters: transport, clientProtocol, connectionToken, connectionData (when using hubs), messageId, groupsToken (if the connection belongs to a group)
Optional parameters: queryString
Sample request:

ws://host/signalr/reconnect?transport=webSockets&clientProtocol=1.4&connectionToken=Aa-
aQA&connectionData=%5B%7B%22Name%22:%22hubConnection%22%7D%5D&messageId=d-3104A0A8-H,0%7CL,0%7CM,2%7CK,0&groupsToken=AQ
Sample response: N/A
*/

// TODO: Abort
// TODO: Ping
