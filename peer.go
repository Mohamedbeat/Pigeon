package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Peer struct {
	Conn     net.Conn
	ClientID string
}

func NewPeer(conn net.Conn) *Peer {
	return &Peer{
		Conn: conn,
	}
}

func (p *Peer) readLoop(s *Server) error {
	defer p.Conn.Close()
	reader := bufio.NewReader(p.Conn)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			return err // connection closed on error
		}
		data := strings.Fields(str)
		if len(data) < 1 {
			continue // ignore empty lines, do not close connection
		}
		err = s.Parser.HandleCommand(data, s, p)
		if err != nil {
			if err == UnknownCommandErr || err == MalformedCommandErr || err == TopicNotFoundErr || err == PeerNotFoundErr {
				p.Conn.Write([]byte(fmt.Sprintf("ERROR: %s\n", err.Error())))
				continue // do not close connection for command errors
			}
			return err // close connection for other errors
		}
		if data[0] == DisconnectCMD {
			return nil // close connection on explicit disconnect
		}
	}
}
