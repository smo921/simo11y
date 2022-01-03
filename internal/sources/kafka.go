package sources

import (
	"ar/internal/types"
	"encoding/json"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func Kafka(done chan string) <-chan types.StructuredMessage {
	out := make(chan types.StructuredMessage)

	// setup kafka client
	broker := "localhost:9092"
	topic := "demo_topic"

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":          "demo",
		"auto.offset.reset": "beginning",
	})

	if err != nil {
		fmt.Printf("Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	consumer.Subscribe(topic, nil)

	go func() {
		defer close(out)
		defer consumer.Close()

		for {
			select {
			case <-done:
				return
			default:
				fmt.Printf("Polling for new message on %s\n", topic)
				ev := consumer.Poll(1000)
				if ev == nil {
					continue
				}
				fmt.Printf("Processing event: %v\n", ev)
				switch e := ev.(type) {
				case *kafka.Message:
					var msg map[string]interface{}
					fmt.Printf("%% Message on %s:\n%s\n",
						e.TopicPartition, string(e.Value))
					err := json.Unmarshal(e.Value, &msg)
					if err != nil {
						fmt.Printf("%% Error: %v\n", e)
					}
					out <- msg
				case kafka.PartitionEOF:
					fmt.Printf("%% Reached %v\n", e)
				case kafka.Error:
					fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
				}
			}
		}
	}()
	return out
}
