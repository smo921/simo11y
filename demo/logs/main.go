package main

import (
	"fmt"

	"ar/internal/consumers"
	"ar/internal/filters"
	logGenerator "ar/internal/generator/logs"
	"ar/internal/transformers"
)

const numMessages = 20

func main() {
	fmt.Println("Starting")
	done := make(chan string)
	defer close(done)

	<-consumers.StructuredMessage(done,
		filters.Take(done, numMessages,
			consumers.Processor(done, transformers.LogHash,
				transformers.StructuredMessage(done,
					logGenerator.SteadyStream(done, 2, logGenerator.Messages(done)),
				),
			),
		),
	)
	fmt.Println("All Done")
}
