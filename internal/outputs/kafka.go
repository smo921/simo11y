package outputs

import (
	"ar/internal/types"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// KafkaConfig stores configuration details
type KafkaConfig struct{}

// Kafka sync reads from the in channel and writes messages to the configured brokers/topics
func Kafka(done chan string, in <-chan types.StructuredMessage) {
	// setup kafka client
	broker := "localhost:9092"
	topic := "demo_topic"

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Producer %v\n", p)
	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)

	for {
		select {
		case <-done:
			fmt.Printf("Kafka sync finished.")
			return
		case msg, ok := <-in:
			if !ok {
				return
			}
			// write message to kafka
			err = p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          msg.Raw(),
				Headers:        []kafka.Header{},
			}, deliveryChan)

			if err != nil {
				fmt.Printf("Error producing message: %s", err)
			} else {
				e := <-deliveryChan
				m := e.(*kafka.Message)

				if m.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
				} else {
					fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
						*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
				}
			}
		}
	}
}
