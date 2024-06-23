package pubsub

import (
	"testing"
)

func TestPubSUb_PublishMessageForAllSubscribers(t *testing.T) {
	ps := NewPubSubClient()
	subscriber_cn := ps.Subscribe(Subscriber1)

	ps.AddPublisher(Publisher1)

	ps.PublishMessageForAllSubscribers("message")

	if msg := <-subscriber_cn; msg != "message" {
		t.Fatalf("Expected 'message', got '%s'", msg)
	}

}
