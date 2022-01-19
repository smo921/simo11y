package outputs

import (
	"fmt"
	"os"
	"strconv"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"simo11y/internal/types"
)

// Kafka sync reads from the in channel and writes messages to the configured brokers/topics
func Kafka(done chan string, config types.KafkaConfig, in <-chan types.StructuredMessage) {
	// setup kafka client
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": config.Broker})

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
			fmt.Println("Kafka sync finished.")
			return
		case msg, ok := <-in:
			if !ok {
				return
			}

			// extract kafka message key
			key, err := kafkaKey(msg, config.KeyField)
			if err != nil {
				fmt.Println("ERROR extracting message key:", err)
				key = nil
			}

			// write message to kafka
			err = p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &config.Topic, Partition: kafka.PartitionAny},
				Key:            key,
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
				}
			}
		}
	}
}

func kafkaKey(m types.StructuredMessage, keyField string) ([]byte, error) {
	var key []byte
	k, err := m.Fetch(keyField)
	if err != nil {
		return nil, err
	}
	switch v := k.(type) {
	case int:
		key = []byte(strconv.Itoa(int(v)))
	case string:
		key = []byte(v)
	case float64:
		key = []byte(strconv.Itoa(int(float64(v))))
	default:
		return nil, fmt.Errorf("error converting keyField '%s'", keyField)
	}

	return key, nil
}
