package pubsub

import (
	"testing"
)

const Subscriber1 = "subscriber1"

func TestPubSub_Subscribe(t *testing.T) {
	ps := NewPubSubClient()
	publisher_ch := ps.AddPublisher(Publisher1)
	ch := ps.Subscribe(Subscriber1)
	if ch == nil {
		t.Fatal("expected channel, got nil")
	}

	if msg := <-publisher_ch; msg != NewSubscriberConnected {
		t.Fatalf("Expected '%s', got '%s'", NewSubscriberConnected, msg)
	}
}

func TestPubSub_Unsubscribe(t *testing.T) {
	ps := NewPubSubClient()

	publisher_cn := ps.AddPublisher(Publisher1)

	ps.Unsubscribe(Subscriber1)

	if msg := <-publisher_cn; msg != NoSubscribersConnected {
		t.Fatalf("Expected '%s', got '%s'", NoSubscribersConnected, msg)
	}

}
