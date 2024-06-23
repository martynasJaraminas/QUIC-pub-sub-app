package pubsub

// Message publishing

func (ps *PubSubClient) PublishMessageForAllSubscribers(msg string) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, ch := range ps.subscribers {
		ch <- msg
	}
}
