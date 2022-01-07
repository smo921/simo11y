package main

import (
	"fmt"

	"simo11y/internal/filters"
	logGenerator "simo11y/internal/generator/logs"
	"simo11y/internal/processors"
	"simo11y/internal/transformers"
)

const numMessages = 20

func main() {
	fmt.Println("Starting")
	done := make(chan string)
	defer close(done)

	messages := filters.Take(done, numMessages,
		processors.StructuredMessage(done, transformers.LogHash,
			transformers.StructuredMessage(done,
				logGenerator.SteadyStream(done, 2, logGenerator.Messages(done)),
			),
		),
	)

	for m := range messages {
		fmt.Println(m)
	}
	fmt.Println("All Done")
}
