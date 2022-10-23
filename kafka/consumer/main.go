package main

import (
	"context"
	"fmt"
	"log"

	kafka "github.com/segmentio/kafka-go"
)

func main() {
	// to produce messages
	topic := "my-topic"
	ctx := context.Background()

	forever := make(chan bool)

	go func() {
		groupID := "my-group-1"
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{"localhost:9092"},
			Topic:    topic,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10mb
			GroupID:  groupID,
		})

		for {
			m, err := r.ReadMessage(ctx)
			if err != nil {
				break
			}
			fmt.Printf("Group:%s read message at offset %d: %s = %s\n", groupID, m.Offset, string(m.Key), string(m.Value))
		}

		if err := r.Close(); err != nil {
			log.Fatal("failed to close reader:", err)
		}
	}()

	go func() {
		groupID := "my-group-2"
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{"localhost:9092"},
			Topic:    topic,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10mb
			GroupID:  groupID,
		})

		for {
			m, err := r.ReadMessage(ctx)
			if err != nil {
				break
			}
			fmt.Printf("Group:%s read message at offset %d: %s = %s\n", groupID, m.Offset, string(m.Key), string(m.Value))
		}

		if err := r.Close(); err != nil {
			log.Fatal("failed to close reader:", err)
		}
	}()

	<-forever
}
