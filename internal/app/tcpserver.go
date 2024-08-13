package app

import (
	"net"

	"go.uber.org/zap"
)

func NewTCPServer(listenAddress string, logger *zap.Logger) *TCPServer {
	return &TCPServer{
		listenAddress: listenAddress,
		quitch:        make(chan struct{}),
		MsgChannel:    make(chan TCPMessage),
		Logger:        logger,
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
			s.Logger.Sugar().Errorf("Error accepting new connetions %s", err)
			continue
		}

		s.Logger.Sugar().Infof("New connection from:%s", conn.RemoteAddr())
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
