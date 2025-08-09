package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
)

const (
	NegotiationPath = "/negotiate"

	// TODO: This may need to be used to track separate messages on the signalr transport:
	// https://github.com/dotnet/aspnetcore/blob/main/src/SignalR/docs/specs/HubProtocol.md#overview
	RecordSeparator = "\x1E"
)

type Client struct {
	url     string
	version uint
	ack     bool
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

	if len(negotiationResponse.AvailableTransports) == 0 {
		return errors.New("no available transports")
	}

	var transport AvailableTransport
	var transportPriority = math.MinInt32
	for _, available := range negotiationResponse.AvailableTransports {
		priority := -1
		switch available.Transport {
		case "LongPolling":
			priority = 0
		case "ServerSentEvents":
			priority = 1
		case "WebSockets":
			priority = 2
		}
		if priority > transportPriority {
			transport = available
			transportPriority = priority
		}

	}

	fmt.Println("Chosen transport:", transport)

	return nil
}

func (c *Client) negotiate(url url.URL) (*NegotiationResponseV1, error) {

	url.Query().Add("version", fmt.Sprint(c.version))
	if c.ack {
		url.Query().Add("useAck", fmt.Sprint(c.ack))
	}

	res, err := http.Post(url.String(), "application/json", nil)
	if err != nil || res.StatusCode != 200 {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	var negotiationResponse NegotiationResponseV1
	err = json.Unmarshal(body, &negotiationResponse)
	if err != nil {
		return nil, err
	}

	return &negotiationResponse, nil
}
