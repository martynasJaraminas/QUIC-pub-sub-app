package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/quic-go/quic-go"
)

func StartPublisherServer(ps *PubSub) {
	publisherAddr := os.Getenv("LOCAL_HOST") + ":" + os.Getenv("PUBLISHER_PORT")
	listener, err := quic.ListenAddr(publisherAddr, GenerateTLSConfig(), nil)
	if err != nil {
		log.Fatalf("Failed to start publisher server: %v", err)
	}

	fmt.Println("Publisher server started on", publisherAddr)
	for {
		session, err := listener.Accept(context.Background())
		if err != nil {
			log.Println("Failed to accept session:", err)
		}
		handlePublisherSession(session, ps)
	}
}

func handlePublisherSession(session quic.Connection, ps *PubSub) {
	log.Println("Accepted session")
	id := session.RemoteAddr().String()
	ch := ps.AddPublisher(id)

	stream, err := session.AcceptStream(context.Background())
	if err != nil {
		log.Println("Failed to accept stream:", err)
		return
	}

	go func() {
		log.Println("handlePublisherStream: Waiting for notifications")
		for status := range ch {
			_, err := stream.Write([]byte(status))
			if err != nil {
				log.Println("Failed to notify publisher:", err)
				return
			}
		}
	}()

	ps.NotifyPublisherAboutSUbscribers(id)

	go handlePublisherStream(stream, ps, id)
}

func handlePublisherStream(stream quic.Stream, ps *PubSub, subscriberId string) {
	log.Println("New stream opened for publisher")
	defer stream.Close()

	buf := make([]byte, 1024)

	for {
		n, err := stream.Read(buf)
		if err != nil {
			log.Println("Failed to read from stream:", err)
			return
		}

		msg := string(buf[:n])
		log.Printf("Publishing message by %s : %s", subscriberId, msg)
		ps.Publish(msg)
	}
}
