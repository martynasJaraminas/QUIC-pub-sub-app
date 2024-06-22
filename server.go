package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/lpernett/godotenv"
	"github.com/quic-go/quic-go"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {

	ps := NewPubSub()

	go startPublisherServer(ps)
	go startSubscriberServer(ps)

	select {} // Block forever
}

func startPublisherServer(ps *PubSub) {
	publisherAddr := os.Getenv("LOCAL_HOST") + ":" + os.Getenv("PUBLISHER_PORT")
	listener, err := quic.ListenAddr(publisherAddr, generateTLSConfig(), nil)
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

func startSubscriberServer(ps *PubSub) {
	subscriberAddr := os.Getenv("LOCAL_HOST") + ":" + os.Getenv("SUBSCRIBER_PORT")
	listener, err := quic.ListenAddr(subscriberAddr, generateTLSConfig(), nil)
	if err != nil {
		log.Fatalf("Failed to start subscriber server: %v", err)
	}

	fmt.Println("Subscriber server started on", subscriberAddr)

	for {
		session, err := listener.Accept(context.Background())
		if err != nil {
			log.Println("Failed to accept session:", err)
			continue
		}
		go handleSubscriberSession(session, ps)
	}
}

func handleSubscriberSession(session quic.Connection, ps *PubSub) {
	fmt.Printf("Subscriber connected from %s\n", session.RemoteAddr().String())

	id := session.RemoteAddr().String()
	ch := ps.Subscribe(id)

	stream, err := session.AcceptStream(context.Background())
	if err != nil {
		log.Println("Failed to accept stream:", err)
		return
	}

	defer ps.Unsubscribe(id)
	handleSubscriberStream(stream, ch, id)

}

func handleSubscriberStream(stream quic.Stream, ch chan string, subscriberId string) {
	log.Println("New stream opened for subscriber")
	defer stream.Close()
	for msg := range ch {
		_, err := stream.Write([]byte(msg))
		// TODO: error is only invoked on timeout: no recent network activity
		// Can i make this fail faster?
		if err != nil {
			log.Println("Failed to write to stream:", err)
			return
		}
		fmt.Printf("Sent message to subscriber %s: %s", subscriberId, msg)
	}
}

func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	// Cert needs to be created, other wise publisher returns "tls: unrecognized name"
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-pubsub"},
	}
}
