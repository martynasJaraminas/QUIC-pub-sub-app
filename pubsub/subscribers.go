package pubsub

import "log"

const (
	NewSubscriberConnected = "New subscriber connected"
	NoSubscribersConnected = "No subscribers connected"
)

// Subscribers methods

func (ps *PubSubClient) Subscribe(id string) chan string {
	// lock ensures that only this goroutine access the map at a time
	ps.mu.Lock()
	defer ps.mu.Unlock()
	// if channel is created with size of 0, unit test hangs
	ch := make(chan string, 1)
	ps.subscribers[id] = ch
	ps.notifyAllPublishers(NewSubscriberConnected)
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
		ps.notifyAllPublishers(NoSubscribersConnected)
	}
}
