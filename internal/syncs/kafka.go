package syncs

import (
	"ar/internal/types"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// KafkaConfig stores configuration details
type KafkaConfig struct{}

// Kafka sync reads from the in channel and writes messages to the configured brokers/topics
func Kafka(done chan string, config KafkaConfig, in <-chan types.StructuredMessage) {
	// setup kafka client
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Producer %v\n", p)
	deliveryChan := make(chan kafka.Event)

	go func() {
		for {
			select {
			case msg := <-in:
				// write message to kafka
			case <-done:
				return
			}
		}
	}()
}
