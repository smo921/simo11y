package main

import (
	"fmt"

	"ar/internal/filters"
	logGenerator "ar/internal/generator/logs"
	"ar/internal/outputs"
	"ar/internal/processors"
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
		processors.StructuredMessage(done, transformers.LogHash,
			transformers.StructuredMessage(done,
				logGenerator.SteadyStream(done, 2, logGenerator.Messages(done)),
			),
		),
	))

	fmt.Println("Done sending messages to Kafka")
	for m := range sources.Kafka(done, broker, topic, "demo") {
		fmt.Println(m)
	}
	fmt.Println("All Done")
}
