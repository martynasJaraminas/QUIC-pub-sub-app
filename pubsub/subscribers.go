package pubsub

import "log"

// Subscribers methods

func (ps *PubSubClient) Subscribe(id string) chan string {
	// lock ensures that only this goroutine access the map at a time
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan string)
	ps.subscribers[id] = ch
	ps.notifyAllPublishers("New subscriber connected")
	return ch
}

func (ps *PubSubClient) Unsubscribe(id string) {

	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ch, ok := ps.subscribers[id]; ok {
		close(ch)
		delete(ps.subscribers, id)
		log.Printf("Publisher %s left", id)
	}

	if len(ps.subscribers) == 0 {
		ps.notifyAllPublishers("No subscribers connected")
	}
}
