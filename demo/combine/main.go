package main

import (
	"fmt"

	"ar/internal/consumers"
	"ar/internal/generator/logs"
	"ar/internal/mixers"
	"ar/internal/transformers"
)

const numMessages = 3

func main() {
	fmt.Println("Starting")
	done := make(chan string)
	defer close(done)

	source1 := transformers.Add(done, "source", 1,
		transformers.StructuredMessage(done,
			logs.SteadyStream(done, 1, logs.Messages(done)),
		),
	)

	source2 := transformers.Add(done, "source", 2,
		transformers.StructuredMessage(done,
			logs.SlowStream(done, logs.Messages(done)),
		),
	)

	combined := mixers.Combine(done, source1, source2)
	<-consumers.Structured(done, combined)
	fmt.Println("All Done")
}