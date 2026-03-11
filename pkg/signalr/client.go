package signalr

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
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
		log:     slog.Default(),
		url:     "http://localhost:8080/signalr",
		version: 1,
		ack:     false,
	}

	for _, opt := range opts {
		opt(client)
	}
	return client
}

func (c *Client) Connect(hubs []Hub, topics []string) (chan Message, error) {

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

	// The F1 SignalR server requires the client to send a protocol declaration
	// over the WebSocket before it will begin streaming subscription data.
	// Without this write the server sends only empty group-state updates and
	// never delivers the snapshot or incremental messages.
	if err = c.transport.Write(append(
		[]byte(`{"protocol":"json","version":1}`),
		RecordSeparator,
	)); err != nil {
		return nil, err
	}

	// Read the server's first frame: {"C":"...","S":1,"M":[]}.
	// S:1 means the connection was accepted.
	msg, err := c.transport.Read()
	if err != nil {
		return nil, err
	}
	var connected ConnectedMessage
	if err := json.Unmarshal(msg.Raw, &connected); err != nil || connected.S != 1 {
		c.log.Warn("unexpected connect message", "response", string(msg.Raw))
	}

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

	// TODO: Fixme, actually get a response.
	return InvocationResponse{}, nil
}

// start sends the classic ASP.NET SignalR /start HTTP request, which informs
// the server that the WebSocket transport is ready. The server begins streaming
// after this call returns successfully.
func (c *Client) start(base url.URL, token string, hubs []Hub) error {
	startURL := *base.JoinPath("start")
	q := startURL.Query()
	q.Set("transport", "webSockets")
	q.Set("clientProtocol", "1.5")
	q.Set("connectionToken", token)

	hubPayload := make([]struct{ Name string }, len(hubs))
	for i, h := range hubs {
		hubPayload[i].Name = h.String()
	}
	hubsJSON, err := json.Marshal(hubPayload)
	if err != nil {
		return err
	}
	q.Set("connectionData", string(hubsJSON))
	startURL.RawQuery = q.Encode()

	c.log.Debug("start request", "url", startURL.String())
	resp, err := http.Get(startURL.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	c.log.Debug("start response", "status", resp.StatusCode, "body", string(body))

	var result struct {
		Response string `json:"Response"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("start response (status %d, body %q): %w", resp.StatusCode, string(body), err)
	}
	if result.Response != "started" {
		c.log.Warn("unexpected start response", "response", result.Response)
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

func (c *Client) read() chan Message {

	ch := make(chan Message)

	go func() {
		for {
			msg, err := c.transport.Read()
			if err != nil {
				c.log.Error("transport read error", "error", err)
				close(ch)
				break
			}
			ch <- msg
		}
	}()

	return ch
}
