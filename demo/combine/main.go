package main

import (
	"fmt"

	"ar/internal/generator/logs"
	"ar/internal/mixers"
	"ar/internal/processors"
	"ar/internal/transformers"
	"ar/internal/types"
)

//const numMessages = 3

func main() {
	fmt.Println("Starting")
	done := make(chan string)
	defer close(done)

	source1 := processors.StructuredMessage(done, cb1,
		transformers.StructuredMessage(done,
			logs.SteadyStream(done, 1, logs.Messages(done)),
		),
	)

	source2 := processors.StructuredMessage(done, cb2,
		transformers.StructuredMessage(done,
			logs.SlowStream(done, logs.Messages(done)),
		),
	)

	combined := mixers.Combine(done, source1, source2)
	for m := range combined {
		fmt.Println(m)
	}
	fmt.Println("All Done")
}

func cb1(m types.StructuredMessage) types.StructuredMessage {
	m["source"] = 1
	return m
}

func cb2(m types.StructuredMessage) types.StructuredMessage {
	m["source"] = 2
	return m
}
