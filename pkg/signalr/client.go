package signalr

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	NegotiationPath = "/negotiate"

	// TODO: This may need to be used to track separate messages on the signalr transport:
	// https://github.com/dotnet/aspnetcore/blob/main/src/SignalR/docs/specs/HubProtocol.md#overview
	RecordSeparator = 0x1E
)

type Client struct {
	log       *slog.Logger
	url       string
	base      string
	version   uint
	ack       bool
	token     string
	transport Transport
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
		log: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})),
		url:     "http://localhost:8080/signalr",
		version: 1,
		ack:     false,
	}

	for _, opt := range opts {
		opt(client)
	}
	return client
}

func (c *Client) Connect(hubs []Hub, topics []string) (chan interface{}, error) {

	uri, err := url.Parse(c.url)
	if err != nil {
		c.log.Error(err.Error())
		return nil, err
	}

	negotiationResponse, err := c.negotiate(*(uri.JoinPath(NegotiationPath)))
	if err != nil {
		c.log.Error(err.Error())
		return nil, err
	}

	c.log.Debug(negotiationResponse.String())

	if !negotiationResponse.TryWebSockets {
		c.log.Error(ErrWebsocketsUnsupported.Error())
		return nil, ErrWebsocketsUnsupported
	}

	c.log.Info("Chosen transport: websocket")

	c.transport = &WebsocketTransport{}
	err = c.transport.Connect("livetiming.formula1.com/signalr", negotiationResponse.ConnectionToken, []Hub{"Streaming"})
	if err != nil {
		c.log.Error(err.Error())
		return nil, err
	}
	// defer conn.Close()

	// Send a HandshakeRequest
	err = c.transport.Handshake()

	// Read server response
	msg, err := c.transport.Read()
	// should be {} + 0x1E if successful
	c.log.Debug(string(msg))

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

	_, err = c.Invoke(InvocationRequest{
		Hub:    "Streaming",
		Method: "Subscribe",
		Arguments: []interface{}{
			[]string{
				"Heartbeat", "CarData.z", "Position.z",
				"ExtrapolatedClock", "TopThree", "RcmSeries",
				"TimingStats", "TimingAppData",
				"WeatherData", "TrackStatus", "DriverList",
				"RaceControlMessages", "SessionInfo",
				"SessionData", "LapCount", "TimingData",
			},
		},
		I: 1,
	})
	if err != nil {
		c.log.Error(err.Error())
		return nil, err
	}

	return c.read(), nil
}

func (c *Client) Invoke(req InvocationRequest) (InvocationResponse, error) {
	// TODO: Invoke Subscription: hub.server.invoke("Subscribe", self.topics)
	//{"H":"chathub","M":"Send","A":["JS Client","Test message"],"I":0, "S":{"customProperty" : "abc"}}
	// Send a HandshakeRequest
	b, _ := json.Marshal(map[string]interface{}{
		"H": req.Hub,
		"M": req.Method,
		"A": req.Arguments,
		"I": req.I,
	})
	// append record separator (0x1E)
	b = append(b, RecordSeparator)
	err := c.transport.Write(b)
	if err != nil {
		return InvocationResponse{}, err
	}
	fmt.Println("Sending", string(b))

	// TODO: Fixme, actually get a response.
	return InvocationResponse{}, nil
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

// TODO: Reconnect
func (c *Client) reconnect() {
	// TODO: implement reconnection
	panic("not implemented")
}

// TODO: Abort
func (c *Client) abort() {
	// TODO: implement reconnection
	panic("not implemented")
}

// TODO: Ping
func (c *Client) ping() {
	// TODO: implement reconnection
	panic("not implemented")
}

func (c *Client) read() chan interface{} {

	ch := make(chan interface{})

	go func() {
		fmt.Println("Reading messages...")
		buf := ""
		for {
			message, err := c.transport.Read()
			if err != nil {
				close(ch)
				break
			}
			fmt.Printf("Read: %s\n", message)
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
