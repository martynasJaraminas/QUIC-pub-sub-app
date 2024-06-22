package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"

	"github.com/quic-go/quic-go"
)

const subscriberAddr = "localhost:5555"

func main() {
	ctx := context.Background()

	session, err := quic.DialAddr(ctx, subscriberAddr, generateTLSConfig(), nil)
	if err != nil {
		log.Fatal(err)
	}

	stream, err := session.OpenStreamSync(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// I need initial message to establish stream connection to server
	// Later can be used to subscribe to topic specific topic
	sendMessage(stream, "Hello, Im a subscriber")

	buf := make([]byte, 1024)
	for {
		n, err := stream.Read(buf)
		if err != nil {
			log.Fatal(err)
		}

		msg := string(buf[:n])
		fmt.Println("Received message:", msg)
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
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-pubsub"},
	}
}
