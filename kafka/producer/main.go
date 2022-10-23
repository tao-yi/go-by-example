package main

import (
	"context"
	"log"

	kafka "github.com/segmentio/kafka-go"
)

func main() {
	// to produce messages
	topic := "my-topic"
	ctx := context.Background()

	w := &kafka.Writer{
		Addr:  kafka.TCP("localhost:9092"),
		Topic: topic,
		// the topic will be created if it is missing
		AllowAutoTopicCreation: true,
		// Balancer:               &kafka.LeastBytes{},
		Balancer:     &kafka.Hash{},
		RequiredAcks: kafka.RequireAll,
	}

	err := w.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte("Key-A"),
			Value: []byte("Hello World!"),
		},
		kafka.Message{
			Key:   []byte("Key-B"),
			Value: []byte("One!"),
		},
		kafka.Message{
			Key:   []byte("Key-C"),
			Value: []byte("Two!"),
		},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
