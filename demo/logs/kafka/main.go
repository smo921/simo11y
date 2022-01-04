package main

import (
	"fmt"

	"ar/internal/consumers"
	"ar/internal/filters"
	logGenerator "ar/internal/generator/logs"
	"ar/internal/outputs"
	"ar/internal/sources"
	"ar/internal/transformers"
)

const numMessages = 20
const broker = "localhost:9092"
const topic = "demo_topic"

func main() {
	fmt.Println("Starting")
	done := make(chan string)
	defer close(done)

	outputs.Kafka(done, broker, topic, filters.Take(done, numMessages,
		consumers.Processor(done, transformers.LogHash,
			transformers.StructuredMessage(done,
				logGenerator.SteadyStream(done, 2, logGenerator.Messages(done)),
			),
		),
	))

	fmt.Println("Done sending messages to Kafka")

	<-consumers.StructuredMessage(done, sources.Kafka(done, broker, topic, "demo"))

	fmt.Println("All Done")
}
