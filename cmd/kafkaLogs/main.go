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

func main() {
	fmt.Println("Starting")
	done := make(chan string)
	defer close(done)

	outputs.Kafka(done, filters.Take(done, numMessages,
		transformers.LogHash(done, "logHash",
			transformers.StructuredMessage(done,
				logGenerator.SteadyStream(done, 2, logGenerator.Messages(done)),
			),
		),
	))

	fmt.Println("Done sending messages to Kafka")

	<-consumers.Structured(done, sources.Kafka(done))

	fmt.Println("All Done")
}
