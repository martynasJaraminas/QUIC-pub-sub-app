package pubsub

// Notifications for publishers

func (ps *PubSubClient) notifyAllPublishers(msg string) {
	for _, ch := range ps.publishers {
		ch <- msg
	}
}

func (ps *PubSubClient) NotifyPublisherAboutSUbscribers(id string) {
	if ch, ok := ps.publishers[id]; ok && len(ps.subscribers) == 0 {
		ch <- "No subscribers connected"
	} else {
		ch <- "Server has active subscribers"

	}
}
