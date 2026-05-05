package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

var rdb *redis.Client

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}
func main() {
	initRedis()
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
		var msg ML1ScoredMessage
		err = json.Unmarshal(m.Value, &msg)
		if err != nil {
			log.Println("JSON Error message: ", err)
			continue
		}
		saveStatusToRedis(msg.CaseID, "be1_orchestrator_start")
		fmt.Println("Case received", msg.CaseID)
		fmt.Println("Score: ", msg.Scores.Score)
	}
}

func saveStatusToRedis(caseID string, status string) {
	ctx := context.Background()
	key := "case:" + caseID + ":status"

	err := rdb.Set(ctx, key, status, 0).Err()
	if err != nil {
		log.Println("Error occur during save to Redis: ", err)
	}
}
