package app

import (
	"encoding/json"
	"net"

	"github.com/elad57/emuq/internal/broker"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type ClientMessageType string

const (
	RESPONSE  ClientMessageType = "Response"
	SUBSCRIBE ClientMessageType = "Subscribe"
	PRODUCE   ClientMessageType = "Produce"
)

type TCPServer struct {
	Logger        *zap.Logger
	listenAddress string
	ln            net.Listener
	quitch        chan struct{}
	MsgChannel    chan TCPMessage
}

type TCPMessage struct {
	Type       ClientMessageType `json:"type"`
	Payload    json.RawMessage   `json:"payload"`
	Queue      string            `json:"queue"`
	Enviorment string            `json:"enviorment"`
	Sender     string            `json:"sender"`
	Connection *net.Conn
}

type HttpServer struct {
	Broker *broker.Broker
	Router *mux.Router
	Logger  *zap.Logger
}
