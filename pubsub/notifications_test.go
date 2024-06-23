package pubsub

import (
	"testing"
)

func TestPubSub_NotifyPublisherAboutSUbscribers_NoSubscribers(t *testing.T) {
	ps := NewPubSubClient()

	ch_pub_1 := ps.AddPublisher(Publisher1)

	ps.NotifyPublisherAboutSUbscribers(Publisher1)

	if msg := <-ch_pub_1; msg != NoSubscribers {
		t.Fatalf("Expected '%s', got '%s'", NoSubscribers, msg)
	}

}
func TestPubSub_NotifyPublisherAboutSUbscribers_HasActive(t *testing.T) {
	ps := NewPubSubClient()
	ps.Subscribe(Subscriber1)

	ch_pub_1 := ps.AddPublisher(Publisher1)

	ps.NotifyPublisherAboutSUbscribers(Publisher1)

	if msg := <-ch_pub_1; msg != HasActiveSubscribers {
		t.Fatalf("Expected '%s', got '%s'", HasActiveSubscribers, msg)
	}

}
