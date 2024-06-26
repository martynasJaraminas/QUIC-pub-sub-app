package main

import (
	"context"
	"log"
	"os"

	"quic-pub-sub-app/pubsub"

	"github.com/quic-go/quic-go"
)

func StartPublisherServer(ps *pubsub.PubSubClient) {
	publisherAddr := os.Getenv("LOCAL_HOST") + ":" + os.Getenv("PUBLISHER_PORT")
	listener, err := quic.ListenAddr(publisherAddr, GenerateTLSConfig(), nil)
	if err != nil {
		log.Fatalf("Failed to start publisher server: %v", err)
	}

	log.Println("Publisher server started on", publisherAddr)
	for {
		session, err := listener.Accept(context.Background())
		if err != nil {
			log.Println("Failed to accept session:", err)
		}
		go handlePublisherSession(session, ps)
	}
}

func handlePublisherSession(session quic.Connection, ps *pubsub.PubSubClient) {
	log.Println("Accepted session")
	id := session.RemoteAddr().String()
	ch := ps.AddPublisher(id)

	context.AfterFunc(session.Context(), func() {
		// Does not detect session drop only time expiration
		log.Printf("Publisher Session %s expired", id)
		ps.RemovePublisher(id)
	})

	// In case 1 connection drops and opens stream
	for {
		stream, err := session.AcceptStream(context.Background())
		if err != nil {
			log.Println("Failed to accept stream:", err)
			return
		}
		go initPublisherNotifications(stream, ch)
		ps.NotifyPublisherAboutSUbscribers(id)

		defer stream.Close()
		go handlePublisherStream(stream, ps, id)
	}

}

func handlePublisherStream(stream quic.Stream, ps *pubsub.PubSubClient, subscriberId string) {
	log.Printf("New stream opened for publisher %s", subscriberId)

	buf := make([]byte, 1024)
	for {
		n, err := stream.Read(buf)
		if err != nil {
			log.Println("Failed to read from stream:", err)
			return
		}

		msg := string(buf[:n])
		log.Printf("Publishing message by %s : %s", subscriberId, msg)
		ps.PublishMessageForAllSubscribers(msg)
	}
}

func initPublisherNotifications(stream quic.Stream, ch chan string) {
	for status := range ch {
		_, err := stream.Write([]byte(status))
		if err != nil {
			log.Println("Failed to notify publisher:", err)
			return
		}
	}
}
