package main

import (
	"sync"
)

type Topic struct {
	Name        string
	Messages    []*Message
	Subscribers []*Peer
	rrIndex     int
	mu          sync.Mutex
}

func (t *Topic) AddSubscriber(peer *Peer) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, sub := range t.Subscribers {
		if sub.ClientID == peer.ClientID {
			return // already subscribed
		}
	}
	t.Subscribers = append(t.Subscribers, peer)
}

func (t *Topic) RemoveSubscriber(peerID string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for i, sub := range t.Subscribers {
		if sub.ClientID == peerID {
			t.Subscribers = append(t.Subscribers[:i], t.Subscribers[i+1:]...)
			return
		}
	}
}

func (t *Topic) FindSubscriber(id string) (*Peer, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, sub := range t.Subscribers {
		if sub.ClientID == id {
			return sub, true
		}
	}
	return nil, false
}

func (t *Topic) GetSubscribersSlice() []*Peer {
	t.mu.Lock()
	defer t.mu.Unlock()
	return append([]*Peer(nil), t.Subscribers...)
}
