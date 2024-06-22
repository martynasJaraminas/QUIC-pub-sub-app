package main

import (
	"sync"
)

type PubSub struct {
	subscribers map[string]chan string
	publishers  map[string]chan string
	mu          sync.RWMutex
}

func NewPubSub() *PubSub {
	return &PubSub{
		subscribers: make(map[string]chan string, 1),
		publishers:  make(map[string]chan string, 1),
	}
}

func (ps *PubSub) Subscribe(id string) chan string {
	// lock ensures that only this goroutine access the map at a time
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan string, 1)
	ps.subscribers[id] = ch
	ps.notifyAllPublishers("New subscriber connected")
	return ch
}

func (ps *PubSub) Unsubscribe(id string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ch, ok := ps.subscribers[id]; ok {
		close(ch)
		delete(ps.subscribers, id)
	}

	if len(ps.subscribers) == 0 {
		ps.notifyAllPublishers("No subscribers connected")
	}
}

func (ps *PubSub) AddPublisher(id string) chan string {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ch := make(chan string, 1)
	ps.publishers[id] = ch
	return ch
}

func (ps *PubSub) Publish(msg string) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, ch := range ps.subscribers {
		ch <- msg
	}
}

func (ps *PubSub) notifyAllPublishers(msg string) {
	for _, ch := range ps.publishers {
		ch <- msg
	}
}

func (ps *PubSub) NotifyPublisherAboutSUbscribers(id string) {
	if ch, ok := ps.publishers[id]; ok && len(ps.subscribers) == 0 {
		ch <- "No subscribers connected"
	} else {
		ch <- "Server has active subscribers"

	}
}
