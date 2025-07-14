package main

import (
	"fmt"
	"log/slog"
	"net"
	"time"
)

const DefaultAddr = ":3333"

type Config struct {
	Addr string
}

type Message struct {
	ID        string
	Payload   []byte
	Timestamp time.Time
}

type Topic struct {
	Name        string
	Messages    []*Message
	Subscribers map[string]*Peer
}

type Server struct {
	Ln        net.Listener
	Topics    map[string]*Topic
	Cfg       Config
	Peers     map[string]*Peer
	addPeerCh chan *Peer
	rmPeerCh  chan *Peer
	quitCh    chan struct{}
}

func NewServer(addr string) *Server {
	cfg := Config{
		Addr: addr,
	}
	if addr == "" {
		cfg.Addr = DefaultAddr
	}

	return &Server{
		Cfg:       cfg,
		Topics:    map[string]*Topic{},
		Peers:     map[string]*Peer{},
		addPeerCh: make(chan *Peer),
		rmPeerCh:  make(chan *Peer),
		quitCh:    make(chan struct{}),
	}
}

func (s *Server) Run() error {
	ln, err := net.Listen("tcp", s.Cfg.Addr)
	if err != nil {
		return err
	}
	s.Ln = ln
	go s.loop()

	slog.Info("Server Started on port:", s.Cfg.Addr)
	return s.AcceptLoop()
}

func (s *Server) loop() {
	for {
		select {
		case <-s.quitCh:
			return
		case peer := <-s.rmPeerCh:
			delete(s.Peers, peer.ClientID)
		case peer := <-s.addPeerCh:
			s.Peers[peer.ClientID] = peer
		}
	}
}

func (s *Server) AcceptLoop() error {
	for {
		conn, err := s.Ln.Accept()
		if err != nil {
			slog.Error("Accept loop err ", "error", err)
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn)
	slog.Info("Peer connected", "remoteAddr", peer.Conn.RemoteAddr())
	go func() {
		if err := peer.readLoop(s); err != nil {
			slog.Error("raed loop conn error: ", "error", err)
		}
	}()
}

// Finds peer or returns PeerNotFoundErr
func (s *Server) FindPeer(id string) (*Peer, error) {
	peer, ok := s.Peers[id]
	if !ok {
		return nil, PeerNotFoundErr
	}
	return peer, nil
}
func (s *Server) FindTopic(name string) (*Topic, bool) {
	t, exists := s.Topics[name]
	if !exists {
		return &Topic{}, false
	}
	return t, true
}

func (s *Server) FindTopicAndCreate(name string) *Topic {
	if _, exists := s.Topics[name]; !exists {
		t := &Topic{
			Name:        name,
			Messages:    []*Message{},
			Subscribers: make(map[string]*Peer),
		}
		s.Topics[t.Name] = t
		return t
	}
	return s.Topics[name]
}

func (t *Topic) FindSubscriber(id string) (*Peer, bool) {
	peer, exists := t.Subscribers[id]
	return peer, exists
}

func (t *Topic) Publish(payload []byte) {
	for _, sub := range t.Subscribers {
		sub.Conn.Write([]byte(fmt.Sprintf("MESSAGE %s %s\n", t.Name, payload)))
	}
}
