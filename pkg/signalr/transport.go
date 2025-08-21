package signalr

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type Hub struct {
	Name string `json:"Name"`
}

// TODO: Make this a proper interface
type Transport interface {
	Connect(url string) error
}

type WebsocketTransport struct {
	URL string
}

func (WebsocketTransport) Connect(host, token string, hubs []Hub) (*websocket.Conn, error) {
	u, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	u.Scheme = "wss"

	q := u.Query()
	q.Set("transport", "webSockets")
	q.Set("connectionToken", token)

	hubsJson, err := json.Marshal(hubs)
	if err != nil {
		return nil, err
	}
	q.Set("connectionData", string(hubsJson))

	u.RawQuery = q.Encode()

	retries := 5
	var conn *websocket.Conn
	var resp *http.Response
	for i := 0; i < retries; i++ {
		log.Printf("connecting to %s...", u.String())
		conn, resp, err = websocket.DefaultDialer.Dial(
			u.String(),
			//http.Header{"User-Agent": []string{"Go-SignalR-Client"}},
			http.Header{},
		)
		if err != nil {
			log.Printf("Failed to dial: %v %v", err, resp.StatusCode)
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}
	if err != nil {
		log.Printf("Failed to connect after retries: %v", err)
		return nil, err
	}
	return conn, nil
}

func (WebsocketTransport) Send(conn *websocket.Conn) {}
