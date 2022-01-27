package main

import (
	"fmt"
	logGenerator "simo11y/internal/generator/logs"
	"simo11y/internal/mixers"
	"simo11y/internal/processors"
	"simo11y/internal/transformers"
)

const numMessages = 20

func main() {
	fmt.Println("Starting")
	done := make(chan string)
	defer close(done)

	messages := processors.Decompressor(done, processors.Compressor(done,
		mixers.Collector(done, numMessages, -1,
			processors.StructuredMessage(done, transformers.LogHash,
				transformers.StructuredMessage(done,
					logGenerator.SteadyStream(done, 5, logGenerator.Messages(done)),
				),
			),
		),
	))

	for m := range messages {
		fmt.Printf("Collection Len: %d\n", len(m))
		fmt.Println("First Message:", m[0])
		fmt.Println("Last Message:", m[len(m)-1])
	}
	fmt.Println("All Done")
}
