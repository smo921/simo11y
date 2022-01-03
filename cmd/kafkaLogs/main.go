package main

import "fmt"

import logGenerator "ar/internal/generator/logs"
import "ar/internal/filters"
import "ar/internal/syncs"
import "ar/internal/transformers"

const numMessages = 20

func main() {
	fmt.Println("Starting")
	done := make(chan string)
	defer close(done)

	syncs.Kafka(done, syncs.KafkaConfig{},
		filters.Take(done, numMessages,
			transformers.LogHash(done, "logHash",
				transformers.StructuredMessage(done,
					logGenerator.SteadyStream(done, 2, logGenerator.Messages(done)),
				),
			),
		),
	)
	fmt.Println("All Done")
}