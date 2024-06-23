package pubsub

import (
	"sync"
)

type PubSubClient struct {
	subscribers map[string]chan string
	publishers  map[string]chan string
	mu          sync.RWMutex
}

func NewPubSub() *PubSubClient {
	// Get the pointer to the PubSub
	return &PubSubClient{
		subscribers: make(map[string]chan string, 0),
		publishers:  make(map[string]chan string, 0),
	}
}
