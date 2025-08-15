package main

import (
	"context"
	"fmt"
	"github.com/go-kit/log"
	"github.com/philippseith/signalr"
	"net"
	"os"
)

type ReplayHub struct {
}

func (r ReplayHub) Initialize(hubContext signalr.HubContext) {
	//TODO implement me
	panic("implement me")
}

func (r ReplayHub) OnConnected(connectionID string) {
	//TODO implement me
	panic("implement me")
}

func (r ReplayHub) OnDisconnected(connectionID string) {
	//TODO implement me
	panic("implement me")
}

func main() {
	ctx := context.Background()

	hub := ReplayHub{}

	// Typical server with log level debug to Stderr
	server, err := signalr.NewServer(ctx,
		signalr.SimpleHubFactory(hub),
		signalr.Logger(log.NewLogfmtLogger(os.Stderr), true),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Listening on localhost:6502")

	// Serving over TCP, accepting client who use MessagePack or JSON
	addr, _ := net.ResolveTCPAddr("tcp", "localhost:6502")
	listener, _ := net.ListenTCP("tcp", addr)
	tcpConn, _ := listener.Accept()

	err = server.Serve(signalr.NewNetConnection(ctx, tcpConn))
	if err != nil {
		panic(err)
	}
}
