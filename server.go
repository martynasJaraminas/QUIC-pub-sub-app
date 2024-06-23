package main

import (
	"log"

	"quic-pub-sub-app/pubsub"

	"github.com/lpernett/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {

	ps := pubsub.NewPubSub()

	go StartPublisherServer(ps)
	go StartSubscriberServer(ps)

	select {} // Block forever
}
