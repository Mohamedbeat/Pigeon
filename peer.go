package main

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"net"
	"strings"
	"time"

	"github.com/google/uuid"
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
			return err
		}
		data := strings.Fields(str)
		if len(data) < 1 {
			return UnknownCommandErr
		}
		command := data[0]

		// switch on command
		switch command {

		case DisconnectCMD:
			if len(data) < 2 {
				return MalformedCommandErr
			}
			peer, err := s.FindPeer(data[1])
			if err != nil {
				return err
			}
			s.rmPeerCh <- peer
			res := fmt.Sprintf("DISCONNECTED\n")
			p.Conn.Write([]byte(res))
			slog.Info("Peer disconnected with ID: %s", "DISCONNECT", p.ClientID)
			return nil

		case ConnectCMD:
			id := uuid.NewString()
			p.ClientID = id
			s.addPeerCh <- p
			msg := fmt.Sprintf("Peer connected with ID: %s\n", p.ClientID)
			slog.Info(msg, "INFO")
			res := fmt.Sprintf("CONNECTED %s\n", id)
			p.Conn.Write([]byte(res))

		case SubscribeCMD:
			if len(data) < 3 || data[2] == "" {
				return MalformedCommandErr
			}
			peer, err := s.FindPeer(data[1])
			if err != nil {
				return err
			}

			topic := s.FindTopicAndCreate(data[2])
			_, exists := topic.FindSubscriber(peer.ClientID)
			if exists {
				res := fmt.Sprintf("SUBSCRIBED\n")
				p.Conn.Write([]byte(res))
			}

			topic.Subscribers[peer.ClientID] = peer

			res := fmt.Sprintf("SUBSCRIBED\n")
			p.Conn.Write([]byte(res))

			msg := fmt.Sprintf("Peer with ID %s subscribed to topic %s\n", p.ClientID, topic.Name)
			slog.Info(msg, "INFO")

		case PublishCMD:
			if len(data) < 4 || data[2] == "" || data[3] == "" {
				return MalformedCommandErr
			}
			_, err := s.FindPeer(data[1])
			if err != nil {
				return err
			}
			t, exists := s.FindTopic(data[2])
			if !exists {
				return TopicNotFoundErr
			}

			payload := strings.Join(data[3:], " ")
			msg := &Message{
				ID:        uuid.NewString(),
				Payload:   []byte(payload),
				Timestamp: time.Now(),
			}
			t.Messages = append(t.Messages, msg)
			t.Publish(msg.Payload)

		case UnsubscribeCMD:
			log.Println("got:", command)
		default:
			return UnknownCommandErr
		}
	}
}
