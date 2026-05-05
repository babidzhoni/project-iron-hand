package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
)

func main() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		GroupID: "be1-orchestrator",
		Topic:   "ml1-scored",
	})
	defer reader.Close()

	fmt.Println("Start consuming...")

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Error message: ", err)
			continue
		}
		fmt.Printf("Received: %s\n", string(m.Value))
	}
}
