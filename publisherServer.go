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
		handlePublisherSession(session, ps)
	}
}

func handlePublisherSession(session quic.Connection, ps *pubsub.PubSubClient) {
	log.Println("Accepted session")
	id := session.RemoteAddr().String()
	ch := ps.AddPublisher(id)

	stream, err := session.AcceptStream(context.Background())
	if err != nil {
		log.Println("Failed to accept stream:", err)
		return
	}

	go initPublisherNotifications(stream, ch)
	ps.NotifyPublisherAboutSUbscribers(id)
	go handlePublisherStream(stream, ps, id)
}

func handlePublisherStream(stream quic.Stream, ps *pubsub.PubSubClient, subscriberId string) {
	log.Println("New stream opened for publisher")
	defer func() {
		stream.Close()
		ps.RemovePublisher(subscriberId)
	}()

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
	log.Println("handlePublisherStream: Waiting for notifications")
	for status := range ch {
		_, err := stream.Write([]byte(status))
		if err != nil {
			log.Println("Failed to notify publisher:", err)
			return
		}
	}
}
