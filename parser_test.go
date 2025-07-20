package main

import (
	"bytes"
	"net"
	"testing"
	"time"
)

type mockConn struct {
	written bytes.Buffer
}

func (m *mockConn) Read(b []byte) (int, error)         { return 0, nil }
func (m *mockConn) Write(b []byte) (int, error)        { return m.written.Write(b) }
func (m *mockConn) Close() error                       { return nil }
func (m *mockConn) LocalAddr() net.Addr                { return nil }
func (m *mockConn) RemoteAddr() net.Addr               { return nil }
func (m *mockConn) SetDeadline(t time.Time) error      { return nil }
func (m *mockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *mockConn) SetWriteDeadline(t time.Time) error { return nil }

func TestHandleDisconnect(t *testing.T) {
	s := &Server{
		rmPeerCh: make(chan *Peer, 1),
	}
	peer := &Peer{Conn: &mockConn{}, ClientID: "test-id"}
	parser := NewParser()

	err := parser.handleDisconnect(s, peer)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check that peer was sent to rmPeerCh
	select {
	case p := <-s.rmPeerCh:
		if p != peer {
			t.Errorf("expected peer to be sent to rmPeerCh")
		}
	default:
		t.Errorf("expected peer to be sent to rmPeerCh, but channel was empty")
	}

	// Check that correct message was written to Conn
	mock := peer.Conn.(*mockConn)
	expected := "DISCONNECTED\n"
	if mock.written.String() != expected {
		t.Errorf("expected '%s', got '%s'", expected, mock.written.String())
	}
}
