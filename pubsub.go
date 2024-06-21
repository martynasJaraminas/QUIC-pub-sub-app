package main

import (
	"log"
	"sync"
)

type PubSub struct {
	subscribers map[string]chan string
	mu          sync.RWMutex
	pubNotify   chan string
}

func NewPubSub() *PubSub {
	return &PubSub{
		subscribers: make(map[string]chan string, 1),
		pubNotify:   make(chan string, 1),
	}
}

func (ps *PubSub) Subscribe(id string) chan string {
	// lock ensures that only this goroutines access the map at a time
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan string, 1)
	ps.subscribers[id] = ch
	ps.notifyPublishers()
	return ch
}

func (ps *PubSub) Unsubscribe(id string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ch, ok := ps.subscribers[id]; ok {
		close(ch)
		delete(ps.subscribers, id)
	}
	ps.notifyPublishers()
}

func (ps *PubSub) Publish(msg string) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, ch := range ps.subscribers {
		log.Println("Publishing message:", msg)
		ch <- msg
	}
}

func (ps *PubSub) notifyPublishers() {
	var status string
	if len(ps.subscribers) == 0 {
		status = "No subscribers connected"
		ps.pubNotify <- status
	} else {
		status = "Subscriber connected"
		ps.pubNotify <- status
	}
}

func (ps *PubSub) PublisherNotifyChannel() <-chan string {
	return ps.pubNotify
}
