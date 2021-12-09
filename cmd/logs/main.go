package main

import "fmt"

import "ar/internal/consumers"
import "ar/internal/generator/logs"
import "ar/internal/filters"
import "ar/internal/transformers"

const numMessages = 20

func main() {
	fmt.Println("Starting")
	done := make(chan string)
	defer close(done)

	<-consumers.Structured(done,
		filters.Take(done, numMessages,
			transformers.LogHash(done, "logHash",
				transformers.StructuredMessage(done,
					logs.SteadyStream(done, 2, logs.LogMessages(done)),
				),
			),
		),
	)
	fmt.Println("All Done")
}
