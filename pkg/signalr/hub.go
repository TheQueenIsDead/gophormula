package signalr

type Hub string

func (h Hub) String() string {
	return string(h)
}
