package broker

import (
	"net"
	"time"
)

type Message struct {
	ID        string
	Body      []byte
	Timestamp time.Time
}

type Queue struct {
	Name        string
	Description string
	CreatedAt   time.Time
	Subscribers []Consumer
	Messages    []Message
}

type Enviorment struct {
	Name   string
	Queues map[string]Queue
}

type BrokerState struct {
	Enviorments map[string]Enviorment
}

type Consumer struct {
	connection *net.Conn
	offset     int
	name       string
}
