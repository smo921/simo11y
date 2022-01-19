package types

// KafkaConfig stores configuration details
type KafkaConfig struct {
	Broker, Topic, KeyField, ConsumerGroup string
}
