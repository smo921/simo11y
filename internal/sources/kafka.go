package sources

import (
	"simo11y/internal/types"
	"encoding/json"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func Kafka(done chan string, broker, topic, consumerGroup string) <-chan types.StructuredMessage {
	out := make(chan types.StructuredMessage)

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":          consumerGroup,
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
				ev := consumer.Poll(0)
				if ev == nil {
					continue
				}
				switch e := ev.(type) {
				case *kafka.Message:
					var msg map[string]interface{}
					err := json.Unmarshal(e.Value, &msg)
					if err != nil {
						fmt.Printf("%% Error: %v\n", e)
					}
					out <- msg
				case kafka.Error:
					fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
				}
			}
		}
	}()
	return out
}
