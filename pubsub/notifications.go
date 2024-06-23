package pubsub

const (
	NoSubscribers        = "No subscribers connected"
	HasActiveSubscribers = "Server has active subscribers"
)

// Notifications for publishers

func (ps *PubSubClient) notifyAllPublishers(msg string) {
	for _, ch := range ps.publishers {
		ch <- msg
	}
}

func (ps *PubSubClient) NotifyPublisherAboutSUbscribers(id string) {
	if ch, ok := ps.publishers[id]; ok && len(ps.subscribers) == 0 {
		ch <- NoSubscribers
	} else {
		ch <- HasActiveSubscribers

	}
}
