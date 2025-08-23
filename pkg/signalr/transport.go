package signalr

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

// TODO: Make this a properly generic interface
type Transport interface {
	Connect(host, token string, hubs []Hub) error
	Handshake() error
	Read() ([]byte, error)
	Write([]byte) error
}

type WebsocketTransport struct {
	URL  string
	conn *websocket.Conn
}

func (t *WebsocketTransport) Connect(host, token string, hubs []Hub) error {
	u, err := url.Parse(host)
	if err != nil {
		return err
	}
	u.Scheme = "wss"

	q := u.Query()
	q.Set("transport", "webSockets")
	q.Set("connectionToken", token)

	hubsJson, err := json.Marshal(hubs)
	if err != nil {
		return err
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
		return err
	}

	t.conn = conn

	return nil
}

func (t *WebsocketTransport) Invoke() {}

func (t *WebsocketTransport) Handshake() error {
	b, _ := json.Marshal(HandshakeRequest{
		Protocol: "json",
		Version:  1,
	})
	// append record separator (0x1E)
	b = append(b, RecordSeparator)
	err := t.conn.WriteMessage(websocket.TextMessage, b)

	return err
}

func (t *WebsocketTransport) Read() ([]byte, error) {
	_, msg, err := t.conn.ReadMessage()
	return msg, err
}

func (t *WebsocketTransport) Write(msg []byte) error {
	return t.conn.WriteMessage(websocket.TextMessage, msg)
}
