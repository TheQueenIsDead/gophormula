package signalr

import (
	"encoding/json"
)

type Message struct {
	c string `json:"C,omitempty"`
	g string `json:"G,omitempty"`
	i string `json:"I,omitempty"`

	R json.RawMessage `json:"R,omitempty"`
	M json.RawMessage `json:"M,omitempty"`

	Raw []byte
}

func (m *Message) Data() json.RawMessage {
	if m.R != nil {
		return m.R
	}
	if m.M != nil {
		return m.M
	}
	return nil
}

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

func (nr NegotiationResponse) String() string {
	str, _ := json.MarshalIndent(nr, "", "  ")
	return string(str)
}

// HandshakeRequest	Client	Sent by the client to agree on the message format.
type HandshakeRequest struct {
	Protocol string `json:"protocol"`
	Version  int    `json:"version"`
}

// Invocation Indicates a request to invoke a particular method (the Target) with provided Arguments on the remote endpoint.
type InvocationRequest struct {
	//I – invocation identifier – allows to match up responses with requests
	//H – the name of the hub
	//M – the name of the method
	//A – arguments (an array, can be empty if the method does not have any parameters)
	//S – state – a dictionary containing additional custom data (optional, currently not supported by the C++ client)
	Hub       Hub           `json:"H"`
	Method    string        `json:"M"`
	Arguments []interface{} `json:"A"`
	I         interface{}   `json:"I"`
}

type InvocationResponse struct {
	//I – invocation Id (always present)
	//R – the value returned by the server method (present if the method is not void)
	//E – error message
	//H – true if this is a hub error
	//D – an object containing additional error data (can only be present for hub errors)
	//T – stack trace (if detailed error reporting (i.e. the HubConfiguration.EnableDetailedErrors property) is turned on on the server). Note that none of the clients currently propagate the stack trace to the user but if tracing is turned on it will be logged with the message
	//S – state – a dictionary containing additional custom data (optional, currently not supported by the C++ client)
}
