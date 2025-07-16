package main

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	ConnectCMD     = "CONNECT"
	SubscribeCMD   = "SUBSCRIBE"
	UnsubscribeCMD = "UNSUBSCRIBE"
	PublishCMD     = "PUBLISH"
	DisconnectCMD  = "DISCONNECT"
	ListCMD        = "LIST"
	PingCMD        = "PING"
)

var UnknownCommandErr = errors.New("unknown command")
var MalformedCommandErr = errors.New("Malformed command")
var PeerNotFoundErr = errors.New("Peer not found")
var TopicNotFoundErr = errors.New("Topic not found")

type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) handleConnect(s *Server, peer *Peer) error {
	id := uuid.NewString()
	peer.ClientID = id
	s.addPeerCh <- peer
	msg := fmt.Sprintf("Peer connected with ID: %s\n", peer.ClientID)
	slog.Info(msg, "INFO")
	res := fmt.Sprintf("CONNECTED %s\n", id)
	peer.Conn.Write([]byte(res))
	return nil
}

func (p *Parser) handleDisconnect(s *Server, peer *Peer) error {
	s.rmPeerCh <- peer
	res := fmt.Sprintf("DISCONNECTED\n")
	peer.Conn.Write([]byte(res))
	slog.Info("Peer disconnected with ID: %s", "DISCONNECT", peer.ClientID)
	return nil
}

// subscribe command == SUBSCRIBE PEERID TOPIC
func (p *Parser) handleSubscribe(data []string, s *Server, peer *Peer) error {
	if len(data) < 3 || data[1] == "" || data[2] == "" {
		return MalformedCommandErr
	}
	topic := s.FindTopicAndCreate(data[2])
	_, exists := topic.FindSubscriber(peer.ClientID)
	if exists {
		res := fmt.Sprintf("SUBSCRIBED\n")
		peer.Conn.Write([]byte(res))
		return nil
	}
	topic.AddSubscriber(peer)
	res := fmt.Sprintf("SUBSCRIBED\n")
	peer.Conn.Write([]byte(res))
	msg := fmt.Sprintf("Peer with ID %s subscribed to topic %s\n", peer.ClientID, topic.Name)
	slog.Info(msg, "INFO")
	return nil
}

// publish command == PUBLISH PEERID TOPIC PAYLOAD
func (p *Parser) handlePublish(data []string, s *Server, peer *Peer) error {
	if len(data) < 4 || data[1] == "" || data[2] == "" || data[3] == "" {
		return MalformedCommandErr
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
	t.mu.Lock()
	t.Messages = append(t.Messages, msg)
	subs := t.Subscribers
	if len(subs) == 0 {
		t.mu.Unlock()
		peer.Conn.Write([]byte("NO_SUBSCRIBERS\n"))
		return nil
	}
	sub := subs[t.rrIndex%len(subs)]
	sub.Conn.Write([]byte(fmt.Sprintf("MESSAGE %s %s\n", t.Name, payload)))
	t.rrIndex = (t.rrIndex + 1) % len(subs)
	t.mu.Unlock()
	peer.Conn.Write([]byte("PUBLISHED\n"))
	return nil
}

// unsubscribe command == UNSUBSCRIBE PEERID TOPIC
func (p *Parser) handleUnsubscribe(data []string, s *Server, peer *Peer) error {
	if len(data) < 3 || data[1] == "" || data[2] == "" {
		return MalformedCommandErr
	}
	topic, exists := s.FindTopic(data[2])
	if !exists {
		return TopicNotFoundErr
	}
	topic.RemoveSubscriber(peer.ClientID)
	peer.Conn.Write([]byte("UNSUBSCRIBED\n"))
	msg := fmt.Sprintf("Peer with ID %s unsubscribed from topic %s\n", peer.ClientID, topic.Name)
	slog.Info(msg, "INFO")
	return nil
}

func (p *Parser) HandleCommand(data []string, s *Server, peer *Peer) error {
	if len(data) < 1 {
		return MalformedCommandErr
	}
	command := data[0]
	switch command {
	case DisconnectCMD:
		return p.handleDisconnect(s, peer)
	case ConnectCMD:
		return p.handleConnect(s, peer)
	case SubscribeCMD:
		return p.handleSubscribe(data, s, peer)
	case PublishCMD:
		return p.handlePublish(data, s, peer)
	case UnsubscribeCMD:
		return p.handleUnsubscribe(data, s, peer)
	default:
		return UnknownCommandErr
	}
}
