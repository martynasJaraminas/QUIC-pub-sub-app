package main

import (
	"context"
	"log"
	"os"

	"quic-pub-sub-app/pubsub"

	"github.com/quic-go/quic-go"
)

func StartSubscriberServer(ps *pubsub.PubSubClient) {
	subscriberAddr := os.Getenv("LOCAL_HOST") + ":" + os.Getenv("SUBSCRIBER_PORT")
	listener, err := quic.ListenAddr(subscriberAddr, GenerateTLSConfig(), nil)
	if err != nil {
		log.Fatalf("Failed to start subscriber server: %v", err)
	}

	log.Println("Subscriber server started on", subscriberAddr)

	for {
		session, err := listener.Accept(context.Background())
		if err != nil {
			log.Println("Failed to accept session:", err)
			continue
		}
		go handleSubscriberSession(session, ps)
	}
}

func handleSubscriberSession(session quic.Connection, ps *pubsub.PubSubClient) {
	log.Printf("Subscriber connected from %s\n", session.RemoteAddr().String())

	id := session.RemoteAddr().String()
	ch := ps.Subscribe(id)

	stream, err := session.AcceptStream(context.Background())
	if err != nil {
		log.Println("Failed to accept stream:", err)
		return
	}

	go handleSessionDrop(session.Context(), session, stream, ps, id)

	defer ps.Unsubscribe(id)
	handleSubscriberStream(stream, ch, id)

}

func handleSubscriberStream(stream quic.Stream, ch chan string, subscriberId string) {
	log.Println("New stream opened for subscriber")
	defer stream.Close()

	for msg := range ch {
		_, err := stream.Write([]byte(msg))
		if err != nil {
			log.Println("Failed to write to stream:", err)
			return
		}
		log.Printf("Sent message to subscriber %s: %s", subscriberId, msg)
	}
}

func handleSessionDrop(ctx context.Context, session quic.Connection, stream quic.Stream, ps *pubsub.PubSubClient, id string) {
	<-ctx.Done()
	log.Printf("Session closed from %s\n", session.RemoteAddr().String())
	stream.Close()
	ps.Unsubscribe(id)
}
