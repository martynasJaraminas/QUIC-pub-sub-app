package main

import (
	"log"

	"github.com/lpernett/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {

	ps := NewPubSub()

	go StartPublisherServer(ps)
	go StartSubscriberServer(ps)

	select {} // Block forever
}
