package pubsub

import "log"

// Publishers methods

func (ps *PubSubClient) AddPublisher(id string) chan string {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ch := make(chan string)
	ps.publishers[id] = ch
	return ch
}

func (ps *PubSubClient) RemovePublisher(id string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ch, ok := ps.publishers[id]; ok {
		close(ch)
		delete(ps.publishers, id)
		log.Printf("Publisher %s left", id)
	}
}
