package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/quic-go/quic-go"
)

const publisherAddr = "localhost:4444"

func main() {
	ctx := context.Background()
	session, err := quic.DialAddr(ctx, publisherAddr, generateTLSConfig(), nil)
	if err != nil {
		log.Fatal(err)
	}

	stream, err := session.OpenStreamSync(ctx)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Println("Listening for notifications from server")
		buf := make([]byte, 1024)
		for {
			n, err := stream.Read(buf)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Notification from server:", string(buf[:n]))
		}
	}()
	// I need initial message to establish stream connection to server
	// Later can be used to publish topic before sending data
	sendMessage(stream, "Initial message from publisher")

	startMessagingFlow(stream, 10)
}

func startMessagingFlow(stream quic.Stream, intervalInSeconds time.Duration) {
	ticker := time.NewTicker(intervalInSeconds * time.Second)
	defer ticker.Stop() // Make sure to stop the ticker when we're done

	var messageCount int
	for {
		select {
		case <-ticker.C:

			message := "Hello, QUIC Pub/Sub! " + fmt.Sprint(messageCount)
			sendMessage(stream, message)
			messageCount++
			fmt.Println("Message sent:", message)

		}
	}
}

func sendMessage(stream quic.Stream, message string) {
	_, err := stream.Write([]byte(message))
	if err != nil {
		log.Fatal(err)
	}
}

func generateTLSConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true,                    // Only for testing; disable in production
		NextProtos:         []string{"quic-pubsub"}, // Ensure this matches server
	}
}
