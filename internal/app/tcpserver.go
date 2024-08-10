package app

import (
	"fmt"
	"net"
)

func NewTCPServer(listenAddress string) *TCPServer {
	return &TCPServer{
		listenAddress: listenAddress,
		quitch:        make(chan struct{}),
		MsgChannel:    make(chan TCPMessage),
	}
}

func (s *TCPServer) Start() error {
	ln, err := net.Listen("tcp", s.listenAddress)

	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln

	go s.accepctConnections()

	<-s.quitch

	return nil
}

func (s *TCPServer) accepctConnections() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			println("Error accepting new connetions", err)
			continue
		}

		fmt.Println("New connection from:", conn.RemoteAddr())
		go s.readMessageInConnetion(conn)
	}
}

func (s *TCPServer) readMessageInConnetion(conn net.Conn) {
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			continue
		}
		s.MsgChannel <- TCPMessage{
			Payload:    buf[:n],
			Connection: &conn,
		}
	}
}
