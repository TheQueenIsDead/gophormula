package client

type AvailableTransport struct {
	Transport       string   `json:"transport"`
	TransferFormats []string `json:"transferFormats"`
}
type NegotiationResponse struct {
	NegotiateVersion    int                  `json:"negotiateVersion"`
	ConnectionId        string               `json:"connectionId"`
	AvailableTransports []AvailableTransport `json:"availableTransports"`
}

type NegotiationResponseV1 struct {
	ConnectionToken     string               `json:"connectionToken"`
	ConnectionId        string               `json:"connectionId"`
	NegotiateVersion    int                  `json:"negotiateVersion"`
	AvailableTransports []AvailableTransport `json:"availableTransports"`
}
