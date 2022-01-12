package main

import (
	"fmt"

	"simo11y/internal/filters"
	logGenerator "simo11y/internal/generator/logs"
	"simo11y/internal/outputs"
	"simo11y/internal/processors"
	"simo11y/internal/sources"
	"simo11y/internal/transformers"
)

const numMessages = 20
const broker = "localhost:9092"
const topic = "demo_topic"

func main() {
	fmt.Println("Starting")
	done := make(chan string)
	defer close(done)

	outputs.Kafka(done, broker, topic, "account.id", filters.Take(done, numMessages,
		processors.StructuredMessage(done, transformers.LogHash,
			transformers.StructuredMessage(done,
				logGenerator.SteadyStream(done, 2, logGenerator.Messages(done)),
			),
		),
	))

	fmt.Println("Done sending messages to Kafka")
	for m := range sources.Kafka(done, broker, topic, "demo") {
		fmt.Println(m["account"])
	}
	fmt.Println("All Done")
}
