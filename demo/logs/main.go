package main

import (
	"fmt"

	"ar/internal/filters"
	logGenerator "ar/internal/generator/logs"
	"ar/internal/processors"
	"ar/internal/transformers"
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
