package pubsub

import (
	"testing"
)

const Publisher1 = "publisher1"

func TestPubSub_AddPublisher(t *testing.T) {
	ps := NewPubSubClient()
	publisher_ch := ps.AddPublisher(Publisher1)
	if publisher_ch == nil {
		t.Fatal("expected channel, got nil")
	}
}

func TestPubSub_RemovePublisher(t *testing.T) {
	ps := NewPubSubClient()
	ps.AddPublisher(Publisher1)
	ps.RemovePublisher(Publisher1)

	if _, ok := ps.publishers[Publisher1]; ok {
		t.Fatalf("Expected publisher removed, got %s", Publisher1)
	}
}
