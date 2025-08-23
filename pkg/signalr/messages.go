package signalr

import (
	"encoding/json"
	"time"
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

/*
Message types as described by: https://github.com/dotnet/aspnetcore/blob/main/src/SignalR/docs/specs/HubProtocol.md#overview
*/

//Message Name	Sender	Description
//HandshakeRequest	Client	Sent by the client to agree on the message format.
//HandshakeResponse	Server	Sent by the server as an acknowledgment of the previous HandshakeRequest message. Contains an error if the handshake failed.
//Close	Callee, Caller	Sent by the server when a connection is closed. Contains an error if the connection was closed because of an error. Sent by the client when it's closing the connection, unlikely to contain an error.
//Invocation	Caller	Indicates a request to invoke a particular method (the Target) with provided Arguments on the remote endpoint.
//StreamInvocation	Caller	Indicates a request to invoke a streaming method (the Target) with provided Arguments on the remote endpoint.
//StreamItem	Callee, Caller	Indicates individual items of streamed response data from a previous StreamInvocation message or streamed uploads from an invocation with streamIds.
//Completion	Callee, Caller	Indicates a previous Invocation or StreamInvocation has completed or a stream in an Invocation or StreamInvocation has completed. Contains an error if the invocation concluded with an error or the result of a non-streaming method invocation. The result will be absent for void methods. In case of streaming invocations no further StreamItem messages will be received.
//CancelInvocation	Caller	Sent by the client to cancel a streaming invocation on the server.
//Ping	Caller, Callee	Sent by either party to check if the connection is active.
//Ack	Caller, Callee	Sent by either party to acknowledge that messages have been received up to the provided sequence ID.
//Sequence	Caller, Callee	Sent by either party as the first message when a connection reconnects. Specifies what sequence ID they will start sending messages starting at. Duplicate messages are possible to receive and should be ignored.

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

// Ex: {"C":"d-B18F34F0-tyT,1","S":1,"M":[]} -- Init message, may arrive after initial data
type ConnectResponse struct {
	// MEssage id, present for all non-keep alive messages
	C string `json:"C"`
	// Indicates the transport was initialised, transport message
	S int `json:"S"`
	// An array of data
	M []interface{} `json:"M"`
	// TODO: G, T, L
	// G – groups token – an encrypted string representing group membership
	// T – if the value is 1 the client should transition into the reconnecting state and try to reconnect to the server (i.e. send the reconnect request). The server is sending a message with this property set to 1 if it is being shut down or restarted. Applies to the longPolling transport only.
	// L – the delay between re-establishing poll connections. Applies to the longPolling transport only. Used only by the JavaScript client. Configurable on the server by setting the IConfigurationManager.LongPollDelay property.
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

// PersistentMessage
type PersistentMessage struct {
	//	C – message id, present for all non-KeepAlive messages
	//
	//	M – an array containing actual data.
	//
	//{"C":"d-9B7A6976-B,2|C,2","M":["Welcome!"]}
	//	S – indicates that the transport was initialized (a.k.a. init message)
	//
	//{"C":"s-0,2CDDE7A|1,23ADE88|2,297B01B|3,3997404|4,33239B5","S":1,"M":[]}
	//	G – groups token – an encrypted string representing group membership
	//
	//{"C":"d-6CD4082D-B,0|C,2|D,0","G":"92OXaCStiSZGy5K83cEEt8aR2ocER=","M":[]}
	//	T – if the value is 1 the client should transition into the reconnecting state and try to reconnect to the server (i.e. send the reconnect request). The server is sending a message with this property set to 1 if it is being shut down or restarted. Applies to the longPolling transport only.
	//
	//	L – the delay between re-establishing poll connections. Applies to the longPolling transport only. Used only by the JavaScript client. Configurable on the server by setting the IConfigurationManager.LongPollDelay property.
	//
	//{"C":"d-E9D15DD8-B,4|C,0|D,0","L":2000,
	//	"M":[{"H":"ChatHub","M":"broadcastMessage","A":["C++","msg"]}]}

	C string        `json:"C"`
	G string        `json:"G"`
	M []interface{} `json:"M"`
}
type MiscMessage struct {
	R struct {
		Heartbeat struct {
			Utc time.Time `json:"Utc"`
			Kf  bool      `json:"_kf"`
		} `json:"Heartbeat"`
		ExtrapolatedClock struct {
			Utc           time.Time `json:"Utc"`
			Remaining     string    `json:"Remaining"`
			Extrapolating bool      `json:"Extrapolating"`
		} `json:"ExtrapolatedClock"`
		TopThree struct {
			Withheld bool `json:"Withheld"`
			Lines    []struct {
				Position        string `json:"Position"`
				ShowPosition    bool   `json:"ShowPosition"`
				RacingNumber    string `json:"RacingNumber"`
				Tla             string `json:"Tla"`
				BroadcastName   string `json:"BroadcastName"`
				FullName        string `json:"FullName"`
				FirstName       string `json:"FirstName"`
				LastName        string `json:"LastName"`
				Reference       string `json:"Reference"`
				Team            string `json:"Team"`
				TeamColour      string `json:"TeamColour"`
				LapTime         string `json:"LapTime"`
				LapState        int    `json:"LapState"`
				DiffToAhead     string `json:"DiffToAhead"`
				DiffToLeader    string `json:"DiffToLeader"`
				OverallFastest  bool   `json:"OverallFastest"`
				PersonalFastest bool   `json:"PersonalFastest"`
			} `json:"Lines"`
		} `json:"TopThree"`
		TimingStats struct {
			Withheld bool `json:"Withheld"`
			Lines    struct {
				Field1 struct {
					Line                int    `json:"Line"`
					RacingNumber        string `json:"RacingNumber"`
					PersonalBestLapTime struct {
						Lap      int    `json:"Lap"`
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"PersonalBestLapTime"`
					BestSectors []struct {
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"BestSectors"`
					BestSpeeds struct {
						I1 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I1"`
						I2 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I2"`
						FL struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"FL"`
						ST struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"ST"`
					} `json:"BestSpeeds"`
				} `json:"1"`
				Field2 struct {
					Line                int    `json:"Line"`
					RacingNumber        string `json:"RacingNumber"`
					PersonalBestLapTime struct {
						Lap      int    `json:"Lap"`
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"PersonalBestLapTime"`
					BestSectors []struct {
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"BestSectors"`
					BestSpeeds struct {
						I1 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I1"`
						I2 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I2"`
						FL struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"FL"`
						ST struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"ST"`
					} `json:"BestSpeeds"`
				} `json:"4"`
				Field3 struct {
					Line                int    `json:"Line"`
					RacingNumber        string `json:"RacingNumber"`
					PersonalBestLapTime struct {
						Lap      int    `json:"Lap"`
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"PersonalBestLapTime"`
					BestSectors []struct {
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"BestSectors"`
					BestSpeeds struct {
						I1 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I1"`
						I2 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I2"`
						FL struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"FL"`
						ST struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"ST"`
					} `json:"BestSpeeds"`
				} `json:"5"`
				Field4 struct {
					Line                int    `json:"Line"`
					RacingNumber        string `json:"RacingNumber"`
					PersonalBestLapTime struct {
						Lap      int    `json:"Lap"`
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"PersonalBestLapTime"`
					BestSectors []struct {
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"BestSectors"`
					BestSpeeds struct {
						I1 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I1"`
						I2 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I2"`
						FL struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"FL"`
						ST struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"ST"`
					} `json:"BestSpeeds"`
				} `json:"6"`
				Field5 struct {
					Line                int    `json:"Line"`
					RacingNumber        string `json:"RacingNumber"`
					PersonalBestLapTime struct {
						Lap      int    `json:"Lap"`
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"PersonalBestLapTime"`
					BestSectors []struct {
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"BestSectors"`
					BestSpeeds struct {
						I1 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I1"`
						I2 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I2"`
						FL struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"FL"`
						ST struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"ST"`
					} `json:"BestSpeeds"`
				} `json:"10"`
				Field6 struct {
					Line                int    `json:"Line"`
					RacingNumber        string `json:"RacingNumber"`
					PersonalBestLapTime struct {
						Lap      int    `json:"Lap"`
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"PersonalBestLapTime"`
					BestSectors []struct {
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"BestSectors"`
					BestSpeeds struct {
						I1 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I1"`
						I2 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I2"`
						FL struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"FL"`
						ST struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"ST"`
					} `json:"BestSpeeds"`
				} `json:"12"`
				Field7 struct {
					Line                int    `json:"Line"`
					RacingNumber        string `json:"RacingNumber"`
					PersonalBestLapTime struct {
						Lap      int    `json:"Lap"`
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"PersonalBestLapTime"`
					BestSectors []struct {
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"BestSectors"`
					BestSpeeds struct {
						I1 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I1"`
						I2 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I2"`
						FL struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"FL"`
						ST struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"ST"`
					} `json:"BestSpeeds"`
				} `json:"14"`
				Field8 struct {
					Line                int    `json:"Line"`
					RacingNumber        string `json:"RacingNumber"`
					PersonalBestLapTime struct {
						Lap      int    `json:"Lap"`
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"PersonalBestLapTime"`
					BestSectors []struct {
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"BestSectors"`
					BestSpeeds struct {
						I1 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I1"`
						I2 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I2"`
						FL struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"FL"`
						ST struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"ST"`
					} `json:"BestSpeeds"`
				} `json:"16"`
				Field9 struct {
					Line                int    `json:"Line"`
					RacingNumber        string `json:"RacingNumber"`
					PersonalBestLapTime struct {
						Lap      int    `json:"Lap"`
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"PersonalBestLapTime"`
					BestSectors []struct {
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"BestSectors"`
					BestSpeeds struct {
						I1 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I1"`
						I2 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I2"`
						FL struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"FL"`
						ST struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"ST"`
					} `json:"BestSpeeds"`
				} `json:"18"`
				Field10 struct {
					Line                int    `json:"Line"`
					RacingNumber        string `json:"RacingNumber"`
					PersonalBestLapTime struct {
						Lap      int    `json:"Lap"`
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"PersonalBestLapTime"`
					BestSectors []struct {
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"BestSectors"`
					BestSpeeds struct {
						I1 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I1"`
						I2 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I2"`
						FL struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"FL"`
						ST struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"ST"`
					} `json:"BestSpeeds"`
				} `json:"22"`
				Field11 struct {
					Line                int    `json:"Line"`
					RacingNumber        string `json:"RacingNumber"`
					PersonalBestLapTime struct {
						Lap      int    `json:"Lap"`
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"PersonalBestLapTime"`
					BestSectors []struct {
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"BestSectors"`
					BestSpeeds struct {
						I1 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I1"`
						I2 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I2"`
						FL struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"FL"`
						ST struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"ST"`
					} `json:"BestSpeeds"`
				} `json:"23"`
				Field12 struct {
					Line                int    `json:"Line"`
					RacingNumber        string `json:"RacingNumber"`
					PersonalBestLapTime struct {
						Lap      int    `json:"Lap"`
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"PersonalBestLapTime"`
					BestSectors []struct {
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"BestSectors"`
					BestSpeeds struct {
						I1 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I1"`
						I2 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I2"`
						FL struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"FL"`
						ST struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"ST"`
					} `json:"BestSpeeds"`
				} `json:"27"`
				Field13 struct {
					Line                int    `json:"Line"`
					RacingNumber        string `json:"RacingNumber"`
					PersonalBestLapTime struct {
						Lap      int    `json:"Lap"`
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"PersonalBestLapTime"`
					BestSectors []struct {
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"BestSectors"`
					BestSpeeds struct {
						I1 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I1"`
						I2 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I2"`
						FL struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"FL"`
						ST struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"ST"`
					} `json:"BestSpeeds"`
				} `json:"30"`
				Field14 struct {
					Line                int    `json:"Line"`
					RacingNumber        string `json:"RacingNumber"`
					PersonalBestLapTime struct {
						Lap      int    `json:"Lap"`
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"PersonalBestLapTime"`
					BestSectors []struct {
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"BestSectors"`
					BestSpeeds struct {
						I1 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I1"`
						I2 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I2"`
						FL struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"FL"`
						ST struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"ST"`
					} `json:"BestSpeeds"`
				} `json:"31"`
				Field15 struct {
					Line                int    `json:"Line"`
					RacingNumber        string `json:"RacingNumber"`
					PersonalBestLapTime struct {
						Lap      int    `json:"Lap"`
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"PersonalBestLapTime"`
					BestSectors []struct {
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"BestSectors"`
					BestSpeeds struct {
						I1 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I1"`
						I2 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I2"`
						FL struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"FL"`
						ST struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"ST"`
					} `json:"BestSpeeds"`
				} `json:"43"`
				Field16 struct {
					Line                int    `json:"Line"`
					RacingNumber        string `json:"RacingNumber"`
					PersonalBestLapTime struct {
						Lap      int    `json:"Lap"`
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"PersonalBestLapTime"`
					BestSectors []struct {
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"BestSectors"`
					BestSpeeds struct {
						I1 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I1"`
						I2 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I2"`
						FL struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"FL"`
						ST struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"ST"`
					} `json:"BestSpeeds"`
				} `json:"44"`
				Field17 struct {
					Line                int    `json:"Line"`
					RacingNumber        string `json:"RacingNumber"`
					PersonalBestLapTime struct {
						Lap      int    `json:"Lap"`
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"PersonalBestLapTime"`
					BestSectors []struct {
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"BestSectors"`
					BestSpeeds struct {
						I1 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I1"`
						I2 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I2"`
						FL struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"FL"`
						ST struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"ST"`
					} `json:"BestSpeeds"`
				} `json:"55"`
				Field18 struct {
					Line                int    `json:"Line"`
					RacingNumber        string `json:"RacingNumber"`
					PersonalBestLapTime struct {
						Lap      int    `json:"Lap"`
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"PersonalBestLapTime"`
					BestSectors []struct {
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"BestSectors"`
					BestSpeeds struct {
						I1 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I1"`
						I2 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I2"`
						FL struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"FL"`
						ST struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"ST"`
					} `json:"BestSpeeds"`
				} `json:"63"`
				Field19 struct {
					Line                int    `json:"Line"`
					RacingNumber        string `json:"RacingNumber"`
					PersonalBestLapTime struct {
						Lap      int    `json:"Lap"`
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"PersonalBestLapTime"`
					BestSectors []struct {
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"BestSectors"`
					BestSpeeds struct {
						I1 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I1"`
						I2 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I2"`
						FL struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"FL"`
						ST struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"ST"`
					} `json:"BestSpeeds"`
				} `json:"81"`
				Field20 struct {
					Line                int    `json:"Line"`
					RacingNumber        string `json:"RacingNumber"`
					PersonalBestLapTime struct {
						Lap      int    `json:"Lap"`
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"PersonalBestLapTime"`
					BestSectors []struct {
						Position int    `json:"Position"`
						Value    string `json:"Value"`
					} `json:"BestSectors"`
					BestSpeeds struct {
						I1 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I1"`
						I2 struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"I2"`
						FL struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"FL"`
						ST struct {
							Position int    `json:"Position"`
							Value    string `json:"Value"`
						} `json:"ST"`
					} `json:"BestSpeeds"`
				} `json:"87"`
			} `json:"Lines"`
			SessionType string `json:"SessionType"`
		} `json:"TimingStats"`
		TimingAppData struct {
			Lines struct {
				Field1 struct {
					RacingNumber string `json:"RacingNumber"`
					Line         int    `json:"Line"`
					GridPos      string `json:"GridPos"`
					Stints       []struct {
						LapTime         string `json:"LapTime"`
						LapNumber       int    `json:"LapNumber"`
						LapFlags        int    `json:"LapFlags"`
						Compound        string `json:"Compound"`
						New             string `json:"New"`
						TyresNotChanged string `json:"TyresNotChanged"`
						TotalLaps       int    `json:"TotalLaps"`
						StartLaps       int    `json:"StartLaps"`
					} `json:"Stints"`
				} `json:"16"`
				Field2 struct {
					RacingNumber string `json:"RacingNumber"`
					Line         int    `json:"Line"`
					GridPos      string `json:"GridPos"`
					Stints       []struct {
						LapTime         string `json:"LapTime"`
						LapNumber       int    `json:"LapNumber"`
						LapFlags        int    `json:"LapFlags"`
						Compound        string `json:"Compound"`
						New             string `json:"New"`
						TyresNotChanged string `json:"TyresNotChanged"`
						TotalLaps       int    `json:"TotalLaps"`
						StartLaps       int    `json:"StartLaps"`
					} `json:"Stints"`
				} `json:"81"`
				Field3 struct {
					RacingNumber string `json:"RacingNumber"`
					Line         int    `json:"Line"`
					GridPos      string `json:"GridPos"`
					Stints       []struct {
						LapTime         string `json:"LapTime"`
						LapNumber       int    `json:"LapNumber"`
						LapFlags        int    `json:"LapFlags"`
						Compound        string `json:"Compound"`
						New             string `json:"New"`
						TyresNotChanged string `json:"TyresNotChanged"`
						TotalLaps       int    `json:"TotalLaps"`
						StartLaps       int    `json:"StartLaps"`
					} `json:"Stints"`
				} `json:"4"`
				Field4 struct {
					RacingNumber string `json:"RacingNumber"`
					Line         int    `json:"Line"`
					GridPos      string `json:"GridPos"`
					Stints       []struct {
						LapTime         string `json:"LapTime"`
						LapNumber       int    `json:"LapNumber"`
						LapFlags        int    `json:"LapFlags"`
						Compound        string `json:"Compound"`
						New             string `json:"New"`
						TyresNotChanged string `json:"TyresNotChanged"`
						TotalLaps       int    `json:"TotalLaps"`
						StartLaps       int    `json:"StartLaps"`
					} `json:"Stints"`
				} `json:"63"`
				Field5 struct {
					RacingNumber string `json:"RacingNumber"`
					Line         int    `json:"Line"`
					GridPos      string `json:"GridPos"`
					Stints       []struct {
						LapTime         string `json:"LapTime"`
						LapNumber       int    `json:"LapNumber"`
						LapFlags        int    `json:"LapFlags"`
						Compound        string `json:"Compound"`
						New             string `json:"New"`
						TyresNotChanged string `json:"TyresNotChanged"`
						TotalLaps       int    `json:"TotalLaps"`
						StartLaps       int    `json:"StartLaps"`
					} `json:"Stints"`
				} `json:"14"`
				Field6 struct {
					RacingNumber string `json:"RacingNumber"`
					Line         int    `json:"Line"`
					GridPos      string `json:"GridPos"`
					Stints       []struct {
						LapTime         string `json:"LapTime"`
						LapNumber       int    `json:"LapNumber"`
						LapFlags        int    `json:"LapFlags"`
						Compound        string `json:"Compound"`
						New             string `json:"New"`
						TyresNotChanged string `json:"TyresNotChanged"`
						TotalLaps       int    `json:"TotalLaps"`
						StartLaps       int    `json:"StartLaps"`
					} `json:"Stints"`
				} `json:"18"`
				Field7 struct {
					RacingNumber string `json:"RacingNumber"`
					Line         int    `json:"Line"`
					GridPos      string `json:"GridPos"`
					Stints       []struct {
						LapTime         string `json:"LapTime"`
						LapNumber       int    `json:"LapNumber"`
						LapFlags        int    `json:"LapFlags"`
						Compound        string `json:"Compound"`
						New             string `json:"New"`
						TyresNotChanged string `json:"TyresNotChanged"`
						TotalLaps       int    `json:"TotalLaps"`
						StartLaps       int    `json:"StartLaps"`
					} `json:"Stints"`
				} `json:"5"`
				Field8 struct {
					RacingNumber string `json:"RacingNumber"`
					Line         int    `json:"Line"`
					GridPos      string `json:"GridPos"`
					Stints       []struct {
						LapTime         string `json:"LapTime"`
						LapNumber       int    `json:"LapNumber"`
						LapFlags        int    `json:"LapFlags"`
						Compound        string `json:"Compound"`
						New             string `json:"New"`
						TyresNotChanged string `json:"TyresNotChanged"`
						TotalLaps       int    `json:"TotalLaps"`
						StartLaps       int    `json:"StartLaps"`
					} `json:"Stints"`
				} `json:"1"`
				Field9 struct {
					RacingNumber string `json:"RacingNumber"`
					Line         int    `json:"Line"`
					GridPos      string `json:"GridPos"`
					Stints       []struct {
						LapTime         string `json:"LapTime"`
						LapNumber       int    `json:"LapNumber"`
						LapFlags        int    `json:"LapFlags"`
						Compound        string `json:"Compound"`
						New             string `json:"New"`
						TyresNotChanged string `json:"TyresNotChanged"`
						TotalLaps       int    `json:"TotalLaps"`
						StartLaps       int    `json:"StartLaps"`
					} `json:"Stints"`
				} `json:"30"`
				Field10 struct {
					RacingNumber string `json:"RacingNumber"`
					Line         int    `json:"Line"`
					GridPos      string `json:"GridPos"`
					Stints       []struct {
						LapTime         string `json:"LapTime"`
						LapNumber       int    `json:"LapNumber"`
						LapFlags        int    `json:"LapFlags"`
						Compound        string `json:"Compound"`
						New             string `json:"New"`
						TyresNotChanged string `json:"TyresNotChanged"`
						TotalLaps       int    `json:"TotalLaps"`
						StartLaps       int    `json:"StartLaps"`
					} `json:"Stints"`
				} `json:"6"`
				Field11 struct {
					RacingNumber string `json:"RacingNumber"`
					Line         int    `json:"Line"`
					GridPos      string `json:"GridPos"`
					Stints       []struct {
						LapTime         string `json:"LapTime"`
						LapNumber       int    `json:"LapNumber"`
						LapFlags        int    `json:"LapFlags"`
						Compound        string `json:"Compound"`
						New             string `json:"New"`
						TyresNotChanged string `json:"TyresNotChanged"`
						TotalLaps       int    `json:"TotalLaps"`
						StartLaps       int    `json:"StartLaps"`
					} `json:"Stints"`
				} `json:"87"`
				Field12 struct {
					RacingNumber string `json:"RacingNumber"`
					Line         int    `json:"Line"`
					GridPos      string `json:"GridPos"`
					Stints       []struct {
						LapTime         string `json:"LapTime"`
						LapNumber       int    `json:"LapNumber"`
						LapFlags        int    `json:"LapFlags"`
						Compound        string `json:"Compound"`
						New             string `json:"New"`
						TyresNotChanged string `json:"TyresNotChanged"`
						TotalLaps       int    `json:"TotalLaps"`
						StartLaps       int    `json:"StartLaps"`
					} `json:"Stints"`
				} `json:"44"`
				Field13 struct {
					RacingNumber string `json:"RacingNumber"`
					Line         int    `json:"Line"`
					GridPos      string `json:"GridPos"`
					Stints       []struct {
						LapTime         string `json:"LapTime"`
						LapNumber       int    `json:"LapNumber"`
						LapFlags        int    `json:"LapFlags"`
						Compound        string `json:"Compound"`
						New             string `json:"New"`
						TyresNotChanged string `json:"TyresNotChanged"`
						TotalLaps       int    `json:"TotalLaps"`
						StartLaps       int    `json:"StartLaps"`
					} `json:"Stints"`
				} `json:"55"`
				Field14 struct {
					RacingNumber string `json:"RacingNumber"`
					Line         int    `json:"Line"`
					GridPos      string `json:"GridPos"`
					Stints       []struct {
						LapTime         string `json:"LapTime"`
						LapNumber       int    `json:"LapNumber"`
						LapFlags        int    `json:"LapFlags"`
						Compound        string `json:"Compound"`
						New             string `json:"New"`
						TyresNotChanged string `json:"TyresNotChanged"`
						TotalLaps       int    `json:"TotalLaps"`
						StartLaps       int    `json:"StartLaps"`
					} `json:"Stints"`
				} `json:"43"`
				Field15 struct {
					RacingNumber string `json:"RacingNumber"`
					Line         int    `json:"Line"`
					GridPos      string `json:"GridPos"`
					Stints       []struct {
						LapTime         string `json:"LapTime"`
						LapNumber       int    `json:"LapNumber"`
						LapFlags        int    `json:"LapFlags"`
						Compound        string `json:"Compound"`
						New             string `json:"New"`
						TyresNotChanged string `json:"TyresNotChanged"`
						TotalLaps       int    `json:"TotalLaps"`
						StartLaps       int    `json:"StartLaps"`
					} `json:"Stints"`
				} `json:"12"`
				Field16 struct {
					RacingNumber string `json:"RacingNumber"`
					Line         int    `json:"Line"`
					GridPos      string `json:"GridPos"`
					Stints       []struct {
						LapTime         string `json:"LapTime"`
						LapNumber       int    `json:"LapNumber"`
						LapFlags        int    `json:"LapFlags"`
						Compound        string `json:"Compound"`
						New             string `json:"New"`
						TyresNotChanged string `json:"TyresNotChanged"`
						TotalLaps       int    `json:"TotalLaps"`
						StartLaps       int    `json:"StartLaps"`
					} `json:"Stints"`
				} `json:"10"`
				Field17 struct {
					RacingNumber string `json:"RacingNumber"`
					Line         int    `json:"Line"`
					GridPos      string `json:"GridPos"`
					Stints       []struct {
						LapTime         string `json:"LapTime"`
						LapNumber       int    `json:"LapNumber"`
						LapFlags        int    `json:"LapFlags"`
						Compound        string `json:"Compound"`
						New             string `json:"New"`
						TyresNotChanged string `json:"TyresNotChanged"`
						TotalLaps       int    `json:"TotalLaps"`
						StartLaps       int    `json:"StartLaps"`
					} `json:"Stints"`
				} `json:"31"`
				Field18 struct {
					RacingNumber string `json:"RacingNumber"`
					Line         int    `json:"Line"`
					GridPos      string `json:"GridPos"`
					Stints       []struct {
						LapTime         string `json:"LapTime"`
						LapNumber       int    `json:"LapNumber"`
						LapFlags        int    `json:"LapFlags"`
						Compound        string `json:"Compound"`
						New             string `json:"New"`
						TyresNotChanged string `json:"TyresNotChanged"`
						TotalLaps       int    `json:"TotalLaps"`
						StartLaps       int    `json:"StartLaps"`
					} `json:"Stints"`
				} `json:"27"`
				Field19 struct {
					RacingNumber string `json:"RacingNumber"`
					Line         int    `json:"Line"`
					GridPos      string `json:"GridPos"`
					Stints       []struct {
						LapTime         string `json:"LapTime"`
						LapNumber       int    `json:"LapNumber"`
						LapFlags        int    `json:"LapFlags"`
						Compound        string `json:"Compound"`
						New             string `json:"New"`
						TyresNotChanged string `json:"TyresNotChanged"`
						TotalLaps       int    `json:"TotalLaps"`
						StartLaps       int    `json:"StartLaps"`
					} `json:"Stints"`
				} `json:"23"`
				Field20 struct {
					RacingNumber string `json:"RacingNumber"`
					Line         int    `json:"Line"`
					GridPos      string `json:"GridPos"`
					Stints       []struct {
						LapTime         string `json:"LapTime"`
						LapNumber       int    `json:"LapNumber"`
						LapFlags        int    `json:"LapFlags"`
						Compound        string `json:"Compound"`
						New             string `json:"New"`
						TyresNotChanged string `json:"TyresNotChanged"`
						TotalLaps       int    `json:"TotalLaps"`
						StartLaps       int    `json:"StartLaps"`
					} `json:"Stints"`
				} `json:"22"`
			} `json:"Lines"`
		} `json:"TimingAppData"`
		WeatherData struct {
			AirTemp       string `json:"AirTemp"`
			Humidity      string `json:"Humidity"`
			Pressure      string `json:"Pressure"`
			Rainfall      string `json:"Rainfall"`
			TrackTemp     string `json:"TrackTemp"`
			WindDirection string `json:"WindDirection"`
			WindSpeed     string `json:"WindSpeed"`
			Kf            bool   `json:"_kf"`
		} `json:"WeatherData"`
		TrackStatus struct {
			Status  string `json:"Status"`
			Message string `json:"Message"`
			Kf      bool   `json:"_kf"`
		} `json:"TrackStatus"`
		DriverList struct {
			Field1 struct {
				RacingNumber  string `json:"RacingNumber"`
				BroadcastName string `json:"BroadcastName"`
				FullName      string `json:"FullName"`
				Tla           string `json:"Tla"`
				Line          int    `json:"Line"`
				TeamName      string `json:"TeamName"`
				TeamColour    string `json:"TeamColour"`
				FirstName     string `json:"FirstName"`
				LastName      string `json:"LastName"`
				Reference     string `json:"Reference"`
				HeadshotUrl   string `json:"HeadshotUrl"`
				PublicIdRight string `json:"PublicIdRight"`
			} `json:"16"`
			Field2 struct {
				RacingNumber  string `json:"RacingNumber"`
				BroadcastName string `json:"BroadcastName"`
				FullName      string `json:"FullName"`
				Tla           string `json:"Tla"`
				Line          int    `json:"Line"`
				TeamName      string `json:"TeamName"`
				TeamColour    string `json:"TeamColour"`
				FirstName     string `json:"FirstName"`
				LastName      string `json:"LastName"`
				Reference     string `json:"Reference"`
				HeadshotUrl   string `json:"HeadshotUrl"`
				PublicIdRight string `json:"PublicIdRight"`
			} `json:"81"`
			Field3 struct {
				RacingNumber  string `json:"RacingNumber"`
				BroadcastName string `json:"BroadcastName"`
				FullName      string `json:"FullName"`
				Tla           string `json:"Tla"`
				Line          int    `json:"Line"`
				TeamName      string `json:"TeamName"`
				TeamColour    string `json:"TeamColour"`
				FirstName     string `json:"FirstName"`
				LastName      string `json:"LastName"`
				Reference     string `json:"Reference"`
				HeadshotUrl   string `json:"HeadshotUrl"`
				PublicIdRight string `json:"PublicIdRight"`
			} `json:"4"`
			Field4 struct {
				RacingNumber  string `json:"RacingNumber"`
				BroadcastName string `json:"BroadcastName"`
				FullName      string `json:"FullName"`
				Tla           string `json:"Tla"`
				Line          int    `json:"Line"`
				TeamName      string `json:"TeamName"`
				TeamColour    string `json:"TeamColour"`
				FirstName     string `json:"FirstName"`
				LastName      string `json:"LastName"`
				Reference     string `json:"Reference"`
				HeadshotUrl   string `json:"HeadshotUrl"`
				PublicIdRight string `json:"PublicIdRight"`
			} `json:"63"`
			Field5 struct {
				RacingNumber  string `json:"RacingNumber"`
				BroadcastName string `json:"BroadcastName"`
				FullName      string `json:"FullName"`
				Tla           string `json:"Tla"`
				Line          int    `json:"Line"`
				TeamName      string `json:"TeamName"`
				TeamColour    string `json:"TeamColour"`
				FirstName     string `json:"FirstName"`
				LastName      string `json:"LastName"`
				Reference     string `json:"Reference"`
				HeadshotUrl   string `json:"HeadshotUrl"`
				PublicIdRight string `json:"PublicIdRight"`
			} `json:"14"`
			Field6 struct {
				RacingNumber  string `json:"RacingNumber"`
				BroadcastName string `json:"BroadcastName"`
				FullName      string `json:"FullName"`
				Tla           string `json:"Tla"`
				Line          int    `json:"Line"`
				TeamName      string `json:"TeamName"`
				TeamColour    string `json:"TeamColour"`
				FirstName     string `json:"FirstName"`
				LastName      string `json:"LastName"`
				Reference     string `json:"Reference"`
				HeadshotUrl   string `json:"HeadshotUrl"`
				PublicIdRight string `json:"PublicIdRight"`
			} `json:"18"`
			Field7 struct {
				RacingNumber  string `json:"RacingNumber"`
				BroadcastName string `json:"BroadcastName"`
				FullName      string `json:"FullName"`
				Tla           string `json:"Tla"`
				Line          int    `json:"Line"`
				TeamName      string `json:"TeamName"`
				TeamColour    string `json:"TeamColour"`
				FirstName     string `json:"FirstName"`
				LastName      string `json:"LastName"`
				Reference     string `json:"Reference"`
				HeadshotUrl   string `json:"HeadshotUrl"`
				PublicIdRight string `json:"PublicIdRight"`
			} `json:"5"`
			Field8 struct {
				RacingNumber  string `json:"RacingNumber"`
				BroadcastName string `json:"BroadcastName"`
				FullName      string `json:"FullName"`
				Tla           string `json:"Tla"`
				Line          int    `json:"Line"`
				TeamName      string `json:"TeamName"`
				TeamColour    string `json:"TeamColour"`
				FirstName     string `json:"FirstName"`
				LastName      string `json:"LastName"`
				Reference     string `json:"Reference"`
				HeadshotUrl   string `json:"HeadshotUrl"`
				PublicIdRight string `json:"PublicIdRight"`
			} `json:"1"`
			Field9 struct {
				RacingNumber  string `json:"RacingNumber"`
				BroadcastName string `json:"BroadcastName"`
				FullName      string `json:"FullName"`
				Tla           string `json:"Tla"`
				Line          int    `json:"Line"`
				TeamName      string `json:"TeamName"`
				TeamColour    string `json:"TeamColour"`
				FirstName     string `json:"FirstName"`
				LastName      string `json:"LastName"`
				Reference     string `json:"Reference"`
				HeadshotUrl   string `json:"HeadshotUrl"`
				PublicIdRight string `json:"PublicIdRight"`
			} `json:"30"`
			Field10 struct {
				RacingNumber  string `json:"RacingNumber"`
				BroadcastName string `json:"BroadcastName"`
				FullName      string `json:"FullName"`
				Tla           string `json:"Tla"`
				Line          int    `json:"Line"`
				TeamName      string `json:"TeamName"`
				TeamColour    string `json:"TeamColour"`
				FirstName     string `json:"FirstName"`
				LastName      string `json:"LastName"`
				Reference     string `json:"Reference"`
				HeadshotUrl   string `json:"HeadshotUrl"`
				PublicIdRight string `json:"PublicIdRight"`
			} `json:"6"`
			Field11 struct {
				RacingNumber  string `json:"RacingNumber"`
				BroadcastName string `json:"BroadcastName"`
				FullName      string `json:"FullName"`
				Tla           string `json:"Tla"`
				Line          int    `json:"Line"`
				TeamName      string `json:"TeamName"`
				TeamColour    string `json:"TeamColour"`
				FirstName     string `json:"FirstName"`
				LastName      string `json:"LastName"`
				Reference     string `json:"Reference"`
				HeadshotUrl   string `json:"HeadshotUrl"`
				PublicIdRight string `json:"PublicIdRight"`
			} `json:"87"`
			Field12 struct {
				RacingNumber  string `json:"RacingNumber"`
				BroadcastName string `json:"BroadcastName"`
				FullName      string `json:"FullName"`
				Tla           string `json:"Tla"`
				Line          int    `json:"Line"`
				TeamName      string `json:"TeamName"`
				TeamColour    string `json:"TeamColour"`
				FirstName     string `json:"FirstName"`
				LastName      string `json:"LastName"`
				Reference     string `json:"Reference"`
				HeadshotUrl   string `json:"HeadshotUrl"`
				PublicIdRight string `json:"PublicIdRight"`
			} `json:"44"`
			Field13 struct {
				RacingNumber  string `json:"RacingNumber"`
				BroadcastName string `json:"BroadcastName"`
				FullName      string `json:"FullName"`
				Tla           string `json:"Tla"`
				Line          int    `json:"Line"`
				TeamName      string `json:"TeamName"`
				TeamColour    string `json:"TeamColour"`
				FirstName     string `json:"FirstName"`
				LastName      string `json:"LastName"`
				Reference     string `json:"Reference"`
				HeadshotUrl   string `json:"HeadshotUrl"`
				PublicIdRight string `json:"PublicIdRight"`
			} `json:"55"`
			Field14 struct {
				RacingNumber  string `json:"RacingNumber"`
				BroadcastName string `json:"BroadcastName"`
				FullName      string `json:"FullName"`
				Tla           string `json:"Tla"`
				Line          int    `json:"Line"`
				TeamName      string `json:"TeamName"`
				TeamColour    string `json:"TeamColour"`
				FirstName     string `json:"FirstName"`
				LastName      string `json:"LastName"`
				Reference     string `json:"Reference"`
				PublicIdRight string `json:"PublicIdRight"`
			} `json:"43"`
			Field15 struct {
				RacingNumber  string `json:"RacingNumber"`
				BroadcastName string `json:"BroadcastName"`
				FullName      string `json:"FullName"`
				Tla           string `json:"Tla"`
				Line          int    `json:"Line"`
				TeamName      string `json:"TeamName"`
				TeamColour    string `json:"TeamColour"`
				FirstName     string `json:"FirstName"`
				LastName      string `json:"LastName"`
				Reference     string `json:"Reference"`
				HeadshotUrl   string `json:"HeadshotUrl"`
				PublicIdRight string `json:"PublicIdRight"`
			} `json:"12"`
			Field16 struct {
				RacingNumber  string `json:"RacingNumber"`
				BroadcastName string `json:"BroadcastName"`
				FullName      string `json:"FullName"`
				Tla           string `json:"Tla"`
				Line          int    `json:"Line"`
				TeamName      string `json:"TeamName"`
				TeamColour    string `json:"TeamColour"`
				FirstName     string `json:"FirstName"`
				LastName      string `json:"LastName"`
				Reference     string `json:"Reference"`
				HeadshotUrl   string `json:"HeadshotUrl"`
				PublicIdRight string `json:"PublicIdRight"`
			} `json:"10"`
			Field17 struct {
				RacingNumber  string `json:"RacingNumber"`
				BroadcastName string `json:"BroadcastName"`
				FullName      string `json:"FullName"`
				Tla           string `json:"Tla"`
				Line          int    `json:"Line"`
				TeamName      string `json:"TeamName"`
				TeamColour    string `json:"TeamColour"`
				FirstName     string `json:"FirstName"`
				LastName      string `json:"LastName"`
				Reference     string `json:"Reference"`
				HeadshotUrl   string `json:"HeadshotUrl"`
				PublicIdRight string `json:"PublicIdRight"`
			} `json:"31"`
			Field18 struct {
				RacingNumber  string `json:"RacingNumber"`
				BroadcastName string `json:"BroadcastName"`
				FullName      string `json:"FullName"`
				Tla           string `json:"Tla"`
				Line          int    `json:"Line"`
				TeamName      string `json:"TeamName"`
				TeamColour    string `json:"TeamColour"`
				FirstName     string `json:"FirstName"`
				LastName      string `json:"LastName"`
				Reference     string `json:"Reference"`
				HeadshotUrl   string `json:"HeadshotUrl"`
				PublicIdRight string `json:"PublicIdRight"`
			} `json:"27"`
			Field19 struct {
				RacingNumber  string `json:"RacingNumber"`
				BroadcastName string `json:"BroadcastName"`
				FullName      string `json:"FullName"`
				Tla           string `json:"Tla"`
				Line          int    `json:"Line"`
				TeamName      string `json:"TeamName"`
				TeamColour    string `json:"TeamColour"`
				FirstName     string `json:"FirstName"`
				LastName      string `json:"LastName"`
				Reference     string `json:"Reference"`
				HeadshotUrl   string `json:"HeadshotUrl"`
				PublicIdRight string `json:"PublicIdRight"`
			} `json:"23"`
			Field20 struct {
				RacingNumber  string `json:"RacingNumber"`
				BroadcastName string `json:"BroadcastName"`
				FullName      string `json:"FullName"`
				Tla           string `json:"Tla"`
				Line          int    `json:"Line"`
				TeamName      string `json:"TeamName"`
				TeamColour    string `json:"TeamColour"`
				FirstName     string `json:"FirstName"`
				LastName      string `json:"LastName"`
				Reference     string `json:"Reference"`
				HeadshotUrl   string `json:"HeadshotUrl"`
				PublicIdRight string `json:"PublicIdRight"`
			} `json:"22"`
		} `json:"DriverList"`
		RaceControlMessages struct {
			Messages []struct {
				Utc          string `json:"Utc"`
				Lap          int    `json:"Lap"`
				Category     string `json:"Category"`
				Flag         string `json:"Flag,omitempty"`
				Scope        string `json:"Scope,omitempty"`
				Message      string `json:"Message"`
				Sector       int    `json:"Sector,omitempty"`
				Status       string `json:"Status,omitempty"`
				RacingNumber string `json:"RacingNumber,omitempty"`
			} `json:"Messages"`
		} `json:"RaceControlMessages"`
		SessionInfo struct {
			Meeting struct {
				Key          int    `json:"Key"`
				Name         string `json:"Name"`
				OfficialName string `json:"OfficialName"`
				Location     string `json:"Location"`
				Number       int    `json:"Number"`
				Country      struct {
					Key  int    `json:"Key"`
					Code string `json:"Code"`
					Name string `json:"Name"`
				} `json:"Country"`
				Circuit struct {
					Key       int    `json:"Key"`
					ShortName string `json:"ShortName"`
				} `json:"Circuit"`
			} `json:"Meeting"`
			SessionStatus string `json:"SessionStatus"`
			ArchiveStatus struct {
				Status string `json:"Status"`
			} `json:"ArchiveStatus"`
			Key       int    `json:"Key"`
			Type      string `json:"Type"`
			Name      string `json:"Name"`
			StartDate string `json:"StartDate"`
			EndDate   string `json:"EndDate"`
			GmtOffset string `json:"GmtOffset"`
			Path      string `json:"Path"`
			Kf        bool   `json:"_kf"`
		} `json:"SessionInfo"`
		SessionData struct {
			Series []struct {
				Utc time.Time `json:"Utc"`
				Lap int       `json:"Lap"`
			} `json:"Series"`
			StatusSeries []struct {
				Utc           time.Time `json:"Utc"`
				TrackStatus   string    `json:"TrackStatus,omitempty"`
				SessionStatus string    `json:"SessionStatus,omitempty"`
			} `json:"StatusSeries"`
		} `json:"SessionData"`
		LapCount struct {
			CurrentLap int `json:"CurrentLap"`
			TotalLaps  int `json:"TotalLaps"`
		} `json:"LapCount"`
		TimingData struct {
			Lines struct {
				Field1 struct {
					GapToLeader             string `json:"GapToLeader"`
					IntervalToPositionAhead struct {
						Value    string `json:"Value"`
						Catching bool   `json:"Catching"`
					} `json:"IntervalToPositionAhead"`
					Line             int    `json:"Line"`
					Position         string `json:"Position"`
					ShowPosition     bool   `json:"ShowPosition"`
					RacingNumber     string `json:"RacingNumber"`
					Retired          bool   `json:"Retired"`
					InPit            bool   `json:"InPit"`
					PitOut           bool   `json:"PitOut"`
					Stopped          bool   `json:"Stopped"`
					Status           int    `json:"Status"`
					NumberOfLaps     int    `json:"NumberOfLaps"`
					NumberOfPitStops int    `json:"NumberOfPitStops"`
					Sectors          []struct {
						Stopped       bool   `json:"Stopped"`
						PreviousValue string `json:"PreviousValue"`
						Segments      []struct {
							Status int `json:"Status"`
						} `json:"Segments"`
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"Sectors"`
					Speeds struct {
						I1 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I1"`
						I2 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I2"`
						FL struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"FL"`
						ST struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"ST"`
					} `json:"Speeds"`
					BestLapTime struct {
						Value string `json:"Value"`
						Lap   int    `json:"Lap"`
					} `json:"BestLapTime"`
					LastLapTime struct {
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"LastLapTime"`
				} `json:"16"`
				Field2 struct {
					GapToLeader             string `json:"GapToLeader"`
					IntervalToPositionAhead struct {
						Value    string `json:"Value"`
						Catching bool   `json:"Catching"`
					} `json:"IntervalToPositionAhead"`
					Line             int    `json:"Line"`
					Position         string `json:"Position"`
					ShowPosition     bool   `json:"ShowPosition"`
					RacingNumber     string `json:"RacingNumber"`
					Retired          bool   `json:"Retired"`
					InPit            bool   `json:"InPit"`
					PitOut           bool   `json:"PitOut"`
					Stopped          bool   `json:"Stopped"`
					Status           int    `json:"Status"`
					NumberOfLaps     int    `json:"NumberOfLaps"`
					NumberOfPitStops int    `json:"NumberOfPitStops"`
					Sectors          []struct {
						Stopped       bool   `json:"Stopped"`
						PreviousValue string `json:"PreviousValue"`
						Segments      []struct {
							Status int `json:"Status"`
						} `json:"Segments"`
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"Sectors"`
					Speeds struct {
						I1 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I1"`
						I2 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I2"`
						FL struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"FL"`
						ST struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"ST"`
					} `json:"Speeds"`
					BestLapTime struct {
						Value string `json:"Value"`
						Lap   int    `json:"Lap"`
					} `json:"BestLapTime"`
					LastLapTime struct {
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"LastLapTime"`
				} `json:"81"`
				Field3 struct {
					GapToLeader             string `json:"GapToLeader"`
					IntervalToPositionAhead struct {
						Value    string `json:"Value"`
						Catching bool   `json:"Catching"`
					} `json:"IntervalToPositionAhead"`
					Line             int    `json:"Line"`
					Position         string `json:"Position"`
					ShowPosition     bool   `json:"ShowPosition"`
					RacingNumber     string `json:"RacingNumber"`
					Retired          bool   `json:"Retired"`
					InPit            bool   `json:"InPit"`
					PitOut           bool   `json:"PitOut"`
					Stopped          bool   `json:"Stopped"`
					Status           int    `json:"Status"`
					NumberOfLaps     int    `json:"NumberOfLaps"`
					NumberOfPitStops int    `json:"NumberOfPitStops"`
					Sectors          []struct {
						Stopped       bool   `json:"Stopped"`
						PreviousValue string `json:"PreviousValue"`
						Segments      []struct {
							Status int `json:"Status"`
						} `json:"Segments"`
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"Sectors"`
					Speeds struct {
						I1 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I1"`
						I2 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I2"`
						FL struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"FL"`
						ST struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"ST"`
					} `json:"Speeds"`
					BestLapTime struct {
						Value string `json:"Value"`
						Lap   int    `json:"Lap"`
					} `json:"BestLapTime"`
					LastLapTime struct {
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"LastLapTime"`
				} `json:"4"`
				Field4 struct {
					GapToLeader             string `json:"GapToLeader"`
					IntervalToPositionAhead struct {
						Value    string `json:"Value"`
						Catching bool   `json:"Catching"`
					} `json:"IntervalToPositionAhead"`
					Line             int    `json:"Line"`
					Position         string `json:"Position"`
					ShowPosition     bool   `json:"ShowPosition"`
					RacingNumber     string `json:"RacingNumber"`
					Retired          bool   `json:"Retired"`
					InPit            bool   `json:"InPit"`
					PitOut           bool   `json:"PitOut"`
					Stopped          bool   `json:"Stopped"`
					Status           int    `json:"Status"`
					NumberOfLaps     int    `json:"NumberOfLaps"`
					NumberOfPitStops int    `json:"NumberOfPitStops"`
					Sectors          []struct {
						Stopped       bool   `json:"Stopped"`
						PreviousValue string `json:"PreviousValue"`
						Segments      []struct {
							Status int `json:"Status"`
						} `json:"Segments"`
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"Sectors"`
					Speeds struct {
						I1 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I1"`
						I2 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I2"`
						FL struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"FL"`
						ST struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"ST"`
					} `json:"Speeds"`
					BestLapTime struct {
						Value string `json:"Value"`
						Lap   int    `json:"Lap"`
					} `json:"BestLapTime"`
					LastLapTime struct {
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"LastLapTime"`
				} `json:"63"`
				Field5 struct {
					GapToLeader             string `json:"GapToLeader"`
					IntervalToPositionAhead struct {
						Value    string `json:"Value"`
						Catching bool   `json:"Catching"`
					} `json:"IntervalToPositionAhead"`
					Line             int    `json:"Line"`
					Position         string `json:"Position"`
					ShowPosition     bool   `json:"ShowPosition"`
					RacingNumber     string `json:"RacingNumber"`
					Retired          bool   `json:"Retired"`
					InPit            bool   `json:"InPit"`
					PitOut           bool   `json:"PitOut"`
					Stopped          bool   `json:"Stopped"`
					Status           int    `json:"Status"`
					NumberOfLaps     int    `json:"NumberOfLaps"`
					NumberOfPitStops int    `json:"NumberOfPitStops"`
					Sectors          []struct {
						Stopped       bool   `json:"Stopped"`
						PreviousValue string `json:"PreviousValue"`
						Segments      []struct {
							Status int `json:"Status"`
						} `json:"Segments"`
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"Sectors"`
					Speeds struct {
						I1 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I1"`
						I2 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I2"`
						FL struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"FL"`
						ST struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"ST"`
					} `json:"Speeds"`
					BestLapTime struct {
						Value string `json:"Value"`
						Lap   int    `json:"Lap"`
					} `json:"BestLapTime"`
					LastLapTime struct {
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"LastLapTime"`
				} `json:"14"`
				Field6 struct {
					GapToLeader             string `json:"GapToLeader"`
					IntervalToPositionAhead struct {
						Value    string `json:"Value"`
						Catching bool   `json:"Catching"`
					} `json:"IntervalToPositionAhead"`
					Line             int    `json:"Line"`
					Position         string `json:"Position"`
					ShowPosition     bool   `json:"ShowPosition"`
					RacingNumber     string `json:"RacingNumber"`
					Retired          bool   `json:"Retired"`
					InPit            bool   `json:"InPit"`
					PitOut           bool   `json:"PitOut"`
					Stopped          bool   `json:"Stopped"`
					Status           int    `json:"Status"`
					NumberOfLaps     int    `json:"NumberOfLaps"`
					NumberOfPitStops int    `json:"NumberOfPitStops"`
					Sectors          []struct {
						Stopped       bool   `json:"Stopped"`
						PreviousValue string `json:"PreviousValue"`
						Segments      []struct {
							Status int `json:"Status"`
						} `json:"Segments"`
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"Sectors"`
					Speeds struct {
						I1 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I1"`
						I2 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I2"`
						FL struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"FL"`
						ST struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"ST"`
					} `json:"Speeds"`
					BestLapTime struct {
						Value string `json:"Value"`
						Lap   int    `json:"Lap"`
					} `json:"BestLapTime"`
					LastLapTime struct {
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"LastLapTime"`
				} `json:"18"`
				Field7 struct {
					GapToLeader             string `json:"GapToLeader"`
					IntervalToPositionAhead struct {
						Value    string `json:"Value"`
						Catching bool   `json:"Catching"`
					} `json:"IntervalToPositionAhead"`
					Line             int    `json:"Line"`
					Position         string `json:"Position"`
					ShowPosition     bool   `json:"ShowPosition"`
					RacingNumber     string `json:"RacingNumber"`
					Retired          bool   `json:"Retired"`
					InPit            bool   `json:"InPit"`
					PitOut           bool   `json:"PitOut"`
					Stopped          bool   `json:"Stopped"`
					Status           int    `json:"Status"`
					NumberOfLaps     int    `json:"NumberOfLaps"`
					NumberOfPitStops int    `json:"NumberOfPitStops"`
					Sectors          []struct {
						Stopped       bool   `json:"Stopped"`
						PreviousValue string `json:"PreviousValue"`
						Segments      []struct {
							Status int `json:"Status"`
						} `json:"Segments"`
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"Sectors"`
					Speeds struct {
						I1 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I1"`
						I2 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I2"`
						FL struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"FL"`
						ST struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"ST"`
					} `json:"Speeds"`
					BestLapTime struct {
						Value string `json:"Value"`
						Lap   int    `json:"Lap"`
					} `json:"BestLapTime"`
					LastLapTime struct {
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"LastLapTime"`
				} `json:"5"`
				Field8 struct {
					GapToLeader             string `json:"GapToLeader"`
					IntervalToPositionAhead struct {
						Value    string `json:"Value"`
						Catching bool   `json:"Catching"`
					} `json:"IntervalToPositionAhead"`
					Line             int    `json:"Line"`
					Position         string `json:"Position"`
					ShowPosition     bool   `json:"ShowPosition"`
					RacingNumber     string `json:"RacingNumber"`
					Retired          bool   `json:"Retired"`
					InPit            bool   `json:"InPit"`
					PitOut           bool   `json:"PitOut"`
					Stopped          bool   `json:"Stopped"`
					Status           int    `json:"Status"`
					NumberOfLaps     int    `json:"NumberOfLaps"`
					NumberOfPitStops int    `json:"NumberOfPitStops"`
					Sectors          []struct {
						Stopped       bool   `json:"Stopped"`
						PreviousValue string `json:"PreviousValue"`
						Segments      []struct {
							Status int `json:"Status"`
						} `json:"Segments"`
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"Sectors"`
					Speeds struct {
						I1 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I1"`
						I2 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I2"`
						FL struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"FL"`
						ST struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"ST"`
					} `json:"Speeds"`
					BestLapTime struct {
						Value string `json:"Value"`
						Lap   int    `json:"Lap"`
					} `json:"BestLapTime"`
					LastLapTime struct {
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"LastLapTime"`
				} `json:"1"`
				Field9 struct {
					GapToLeader             string `json:"GapToLeader"`
					IntervalToPositionAhead struct {
						Value    string `json:"Value"`
						Catching bool   `json:"Catching"`
					} `json:"IntervalToPositionAhead"`
					Line             int    `json:"Line"`
					Position         string `json:"Position"`
					ShowPosition     bool   `json:"ShowPosition"`
					RacingNumber     string `json:"RacingNumber"`
					Retired          bool   `json:"Retired"`
					InPit            bool   `json:"InPit"`
					PitOut           bool   `json:"PitOut"`
					Stopped          bool   `json:"Stopped"`
					Status           int    `json:"Status"`
					NumberOfLaps     int    `json:"NumberOfLaps"`
					NumberOfPitStops int    `json:"NumberOfPitStops"`
					Sectors          []struct {
						Stopped       bool   `json:"Stopped"`
						PreviousValue string `json:"PreviousValue"`
						Segments      []struct {
							Status int `json:"Status"`
						} `json:"Segments"`
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"Sectors"`
					Speeds struct {
						I1 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I1"`
						I2 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I2"`
						FL struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"FL"`
						ST struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"ST"`
					} `json:"Speeds"`
					BestLapTime struct {
						Value string `json:"Value"`
						Lap   int    `json:"Lap"`
					} `json:"BestLapTime"`
					LastLapTime struct {
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"LastLapTime"`
				} `json:"30"`
				Field10 struct {
					GapToLeader             string `json:"GapToLeader"`
					IntervalToPositionAhead struct {
						Value    string `json:"Value"`
						Catching bool   `json:"Catching"`
					} `json:"IntervalToPositionAhead"`
					Line             int    `json:"Line"`
					Position         string `json:"Position"`
					ShowPosition     bool   `json:"ShowPosition"`
					RacingNumber     string `json:"RacingNumber"`
					Retired          bool   `json:"Retired"`
					InPit            bool   `json:"InPit"`
					PitOut           bool   `json:"PitOut"`
					Stopped          bool   `json:"Stopped"`
					Status           int    `json:"Status"`
					NumberOfLaps     int    `json:"NumberOfLaps"`
					NumberOfPitStops int    `json:"NumberOfPitStops"`
					Sectors          []struct {
						Stopped       bool   `json:"Stopped"`
						PreviousValue string `json:"PreviousValue"`
						Segments      []struct {
							Status int `json:"Status"`
						} `json:"Segments"`
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"Sectors"`
					Speeds struct {
						I1 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I1"`
						I2 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I2"`
						FL struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"FL"`
						ST struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"ST"`
					} `json:"Speeds"`
					BestLapTime struct {
						Value string `json:"Value"`
						Lap   int    `json:"Lap"`
					} `json:"BestLapTime"`
					LastLapTime struct {
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"LastLapTime"`
				} `json:"6"`
				Field11 struct {
					GapToLeader             string `json:"GapToLeader"`
					IntervalToPositionAhead struct {
						Value    string `json:"Value"`
						Catching bool   `json:"Catching"`
					} `json:"IntervalToPositionAhead"`
					Line             int    `json:"Line"`
					Position         string `json:"Position"`
					ShowPosition     bool   `json:"ShowPosition"`
					RacingNumber     string `json:"RacingNumber"`
					Retired          bool   `json:"Retired"`
					InPit            bool   `json:"InPit"`
					PitOut           bool   `json:"PitOut"`
					Stopped          bool   `json:"Stopped"`
					Status           int    `json:"Status"`
					NumberOfLaps     int    `json:"NumberOfLaps"`
					NumberOfPitStops int    `json:"NumberOfPitStops"`
					Sectors          []struct {
						Stopped       bool   `json:"Stopped"`
						PreviousValue string `json:"PreviousValue"`
						Segments      []struct {
							Status int `json:"Status"`
						} `json:"Segments"`
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"Sectors"`
					Speeds struct {
						I1 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I1"`
						I2 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I2"`
						FL struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"FL"`
						ST struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"ST"`
					} `json:"Speeds"`
					BestLapTime struct {
						Value string `json:"Value"`
						Lap   int    `json:"Lap"`
					} `json:"BestLapTime"`
					LastLapTime struct {
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"LastLapTime"`
				} `json:"87"`
				Field12 struct {
					GapToLeader             string `json:"GapToLeader"`
					IntervalToPositionAhead struct {
						Value    string `json:"Value"`
						Catching bool   `json:"Catching"`
					} `json:"IntervalToPositionAhead"`
					Line             int    `json:"Line"`
					Position         string `json:"Position"`
					ShowPosition     bool   `json:"ShowPosition"`
					RacingNumber     string `json:"RacingNumber"`
					Retired          bool   `json:"Retired"`
					InPit            bool   `json:"InPit"`
					PitOut           bool   `json:"PitOut"`
					Stopped          bool   `json:"Stopped"`
					Status           int    `json:"Status"`
					NumberOfLaps     int    `json:"NumberOfLaps"`
					NumberOfPitStops int    `json:"NumberOfPitStops"`
					Sectors          []struct {
						Stopped       bool   `json:"Stopped"`
						PreviousValue string `json:"PreviousValue"`
						Segments      []struct {
							Status int `json:"Status"`
						} `json:"Segments"`
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"Sectors"`
					Speeds struct {
						I1 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I1"`
						I2 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I2"`
						FL struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"FL"`
						ST struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"ST"`
					} `json:"Speeds"`
					BestLapTime struct {
						Value string `json:"Value"`
						Lap   int    `json:"Lap"`
					} `json:"BestLapTime"`
					LastLapTime struct {
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"LastLapTime"`
				} `json:"44"`
				Field13 struct {
					GapToLeader             string `json:"GapToLeader"`
					IntervalToPositionAhead struct {
						Value    string `json:"Value"`
						Catching bool   `json:"Catching"`
					} `json:"IntervalToPositionAhead"`
					Line             int    `json:"Line"`
					Position         string `json:"Position"`
					ShowPosition     bool   `json:"ShowPosition"`
					RacingNumber     string `json:"RacingNumber"`
					Retired          bool   `json:"Retired"`
					InPit            bool   `json:"InPit"`
					PitOut           bool   `json:"PitOut"`
					Stopped          bool   `json:"Stopped"`
					Status           int    `json:"Status"`
					NumberOfLaps     int    `json:"NumberOfLaps"`
					NumberOfPitStops int    `json:"NumberOfPitStops"`
					Sectors          []struct {
						Stopped       bool   `json:"Stopped"`
						PreviousValue string `json:"PreviousValue"`
						Segments      []struct {
							Status int `json:"Status"`
						} `json:"Segments"`
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"Sectors"`
					Speeds struct {
						I1 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I1"`
						I2 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I2"`
						FL struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"FL"`
						ST struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"ST"`
					} `json:"Speeds"`
					BestLapTime struct {
						Value string `json:"Value"`
						Lap   int    `json:"Lap"`
					} `json:"BestLapTime"`
					LastLapTime struct {
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"LastLapTime"`
				} `json:"55"`
				Field14 struct {
					GapToLeader             string `json:"GapToLeader"`
					IntervalToPositionAhead struct {
						Value    string `json:"Value"`
						Catching bool   `json:"Catching"`
					} `json:"IntervalToPositionAhead"`
					Line             int    `json:"Line"`
					Position         string `json:"Position"`
					ShowPosition     bool   `json:"ShowPosition"`
					RacingNumber     string `json:"RacingNumber"`
					Retired          bool   `json:"Retired"`
					InPit            bool   `json:"InPit"`
					PitOut           bool   `json:"PitOut"`
					Stopped          bool   `json:"Stopped"`
					Status           int    `json:"Status"`
					NumberOfLaps     int    `json:"NumberOfLaps"`
					NumberOfPitStops int    `json:"NumberOfPitStops"`
					Sectors          []struct {
						Stopped       bool   `json:"Stopped"`
						PreviousValue string `json:"PreviousValue"`
						Segments      []struct {
							Status int `json:"Status"`
						} `json:"Segments"`
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"Sectors"`
					Speeds struct {
						I1 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I1"`
						I2 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I2"`
						FL struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"FL"`
						ST struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"ST"`
					} `json:"Speeds"`
					BestLapTime struct {
						Value string `json:"Value"`
						Lap   int    `json:"Lap"`
					} `json:"BestLapTime"`
					LastLapTime struct {
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"LastLapTime"`
				} `json:"43"`
				Field15 struct {
					GapToLeader             string `json:"GapToLeader"`
					IntervalToPositionAhead struct {
						Value    string `json:"Value"`
						Catching bool   `json:"Catching"`
					} `json:"IntervalToPositionAhead"`
					Line             int    `json:"Line"`
					Position         string `json:"Position"`
					ShowPosition     bool   `json:"ShowPosition"`
					RacingNumber     string `json:"RacingNumber"`
					Retired          bool   `json:"Retired"`
					InPit            bool   `json:"InPit"`
					PitOut           bool   `json:"PitOut"`
					Stopped          bool   `json:"Stopped"`
					Status           int    `json:"Status"`
					NumberOfLaps     int    `json:"NumberOfLaps"`
					NumberOfPitStops int    `json:"NumberOfPitStops"`
					Sectors          []struct {
						Stopped       bool   `json:"Stopped"`
						PreviousValue string `json:"PreviousValue"`
						Segments      []struct {
							Status int `json:"Status"`
						} `json:"Segments"`
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"Sectors"`
					Speeds struct {
						I1 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I1"`
						I2 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I2"`
						FL struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"FL"`
						ST struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"ST"`
					} `json:"Speeds"`
					BestLapTime struct {
						Value string `json:"Value"`
						Lap   int    `json:"Lap"`
					} `json:"BestLapTime"`
					LastLapTime struct {
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"LastLapTime"`
				} `json:"12"`
				Field16 struct {
					GapToLeader             string `json:"GapToLeader"`
					IntervalToPositionAhead struct {
						Value    string `json:"Value"`
						Catching bool   `json:"Catching"`
					} `json:"IntervalToPositionAhead"`
					Line             int    `json:"Line"`
					Position         string `json:"Position"`
					ShowPosition     bool   `json:"ShowPosition"`
					RacingNumber     string `json:"RacingNumber"`
					Retired          bool   `json:"Retired"`
					InPit            bool   `json:"InPit"`
					PitOut           bool   `json:"PitOut"`
					Stopped          bool   `json:"Stopped"`
					Status           int    `json:"Status"`
					NumberOfLaps     int    `json:"NumberOfLaps"`
					NumberOfPitStops int    `json:"NumberOfPitStops"`
					Sectors          []struct {
						Stopped       bool   `json:"Stopped"`
						PreviousValue string `json:"PreviousValue"`
						Segments      []struct {
							Status int `json:"Status"`
						} `json:"Segments"`
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"Sectors"`
					Speeds struct {
						I1 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I1"`
						I2 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I2"`
						FL struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"FL"`
						ST struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"ST"`
					} `json:"Speeds"`
					BestLapTime struct {
						Value string `json:"Value"`
						Lap   int    `json:"Lap"`
					} `json:"BestLapTime"`
					LastLapTime struct {
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"LastLapTime"`
				} `json:"10"`
				Field17 struct {
					GapToLeader             string `json:"GapToLeader"`
					IntervalToPositionAhead struct {
						Value    string `json:"Value"`
						Catching bool   `json:"Catching"`
					} `json:"IntervalToPositionAhead"`
					Line             int    `json:"Line"`
					Position         string `json:"Position"`
					ShowPosition     bool   `json:"ShowPosition"`
					RacingNumber     string `json:"RacingNumber"`
					Retired          bool   `json:"Retired"`
					InPit            bool   `json:"InPit"`
					PitOut           bool   `json:"PitOut"`
					Stopped          bool   `json:"Stopped"`
					Status           int    `json:"Status"`
					NumberOfLaps     int    `json:"NumberOfLaps"`
					NumberOfPitStops int    `json:"NumberOfPitStops"`
					Sectors          []struct {
						Stopped       bool   `json:"Stopped"`
						PreviousValue string `json:"PreviousValue"`
						Segments      []struct {
							Status int `json:"Status"`
						} `json:"Segments"`
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"Sectors"`
					Speeds struct {
						I1 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I1"`
						I2 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I2"`
						FL struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"FL"`
						ST struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"ST"`
					} `json:"Speeds"`
					BestLapTime struct {
						Value string `json:"Value"`
						Lap   int    `json:"Lap"`
					} `json:"BestLapTime"`
					LastLapTime struct {
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"LastLapTime"`
				} `json:"31"`
				Field18 struct {
					GapToLeader             string `json:"GapToLeader"`
					IntervalToPositionAhead struct {
						Value    string `json:"Value"`
						Catching bool   `json:"Catching"`
					} `json:"IntervalToPositionAhead"`
					Line             int    `json:"Line"`
					Position         string `json:"Position"`
					ShowPosition     bool   `json:"ShowPosition"`
					RacingNumber     string `json:"RacingNumber"`
					Retired          bool   `json:"Retired"`
					InPit            bool   `json:"InPit"`
					PitOut           bool   `json:"PitOut"`
					Stopped          bool   `json:"Stopped"`
					Status           int    `json:"Status"`
					NumberOfLaps     int    `json:"NumberOfLaps"`
					NumberOfPitStops int    `json:"NumberOfPitStops"`
					Sectors          []struct {
						Stopped       bool   `json:"Stopped"`
						PreviousValue string `json:"PreviousValue"`
						Segments      []struct {
							Status int `json:"Status"`
						} `json:"Segments"`
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"Sectors"`
					Speeds struct {
						I1 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I1"`
						I2 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I2"`
						FL struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"FL"`
						ST struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"ST"`
					} `json:"Speeds"`
					BestLapTime struct {
						Value string `json:"Value"`
						Lap   int    `json:"Lap"`
					} `json:"BestLapTime"`
					LastLapTime struct {
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"LastLapTime"`
				} `json:"27"`
				Field19 struct {
					GapToLeader             string `json:"GapToLeader"`
					IntervalToPositionAhead struct {
						Value    string `json:"Value"`
						Catching bool   `json:"Catching"`
					} `json:"IntervalToPositionAhead"`
					Line             int    `json:"Line"`
					Position         string `json:"Position"`
					ShowPosition     bool   `json:"ShowPosition"`
					RacingNumber     string `json:"RacingNumber"`
					Retired          bool   `json:"Retired"`
					InPit            bool   `json:"InPit"`
					PitOut           bool   `json:"PitOut"`
					Stopped          bool   `json:"Stopped"`
					Status           int    `json:"Status"`
					NumberOfLaps     int    `json:"NumberOfLaps"`
					NumberOfPitStops int    `json:"NumberOfPitStops"`
					Sectors          []struct {
						Stopped       bool   `json:"Stopped"`
						PreviousValue string `json:"PreviousValue"`
						Segments      []struct {
							Status int `json:"Status"`
						} `json:"Segments"`
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"Sectors"`
					Speeds struct {
						I1 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I1"`
						I2 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I2"`
						FL struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"FL"`
						ST struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"ST"`
					} `json:"Speeds"`
					BestLapTime struct {
						Value string `json:"Value"`
						Lap   int    `json:"Lap"`
					} `json:"BestLapTime"`
					LastLapTime struct {
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"LastLapTime"`
				} `json:"23"`
				Field20 struct {
					GapToLeader             string `json:"GapToLeader"`
					IntervalToPositionAhead struct {
						Value    string `json:"Value"`
						Catching bool   `json:"Catching"`
					} `json:"IntervalToPositionAhead"`
					Line             int    `json:"Line"`
					Position         string `json:"Position"`
					ShowPosition     bool   `json:"ShowPosition"`
					RacingNumber     string `json:"RacingNumber"`
					Retired          bool   `json:"Retired"`
					InPit            bool   `json:"InPit"`
					PitOut           bool   `json:"PitOut"`
					Stopped          bool   `json:"Stopped"`
					Status           int    `json:"Status"`
					NumberOfLaps     int    `json:"NumberOfLaps"`
					NumberOfPitStops int    `json:"NumberOfPitStops"`
					Sectors          []struct {
						Stopped       bool   `json:"Stopped"`
						PreviousValue string `json:"PreviousValue"`
						Segments      []struct {
							Status int `json:"Status"`
						} `json:"Segments"`
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"Sectors"`
					Speeds struct {
						I1 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I1"`
						I2 struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"I2"`
						FL struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"FL"`
						ST struct {
							Value           string `json:"Value"`
							Status          int    `json:"Status"`
							OverallFastest  bool   `json:"OverallFastest"`
							PersonalFastest bool   `json:"PersonalFastest"`
						} `json:"ST"`
					} `json:"Speeds"`
					BestLapTime struct {
						Value string `json:"Value"`
						Lap   int    `json:"Lap"`
					} `json:"BestLapTime"`
					LastLapTime struct {
						Value           string `json:"Value"`
						Status          int    `json:"Status"`
						OverallFastest  bool   `json:"OverallFastest"`
						PersonalFastest bool   `json:"PersonalFastest"`
					} `json:"LastLapTime"`
				} `json:"22"`
			} `json:"Lines"`
			Withheld bool `json:"Withheld"`
		} `json:"TimingData"`
	} `json:"R"`
	I string `json:"I"`
}
