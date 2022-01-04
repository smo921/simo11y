package main

import (
	"fmt"

	"ar/internal/consumers"
	"ar/internal/filters"
	"ar/internal/generator/logs"
	"ar/internal/mixers"
	"ar/internal/transformers"
)

const numMessages = 3

func main() {
	fmt.Println("Starting")
	done := make(chan string)
	defer close(done)
	source := transformers.StructuredMessage(done, logs.SteadyStream(done, 1, logs.Messages(done)))
	source = filters.Take(done, numMessages, source)
	ch1, ch2 := mixers.Tee(done, source)
	c1 := consumers.Structured(done, ch1)
	c2 := consumers.Structured(done, ch2)
	<-c1
	<-c2
	fmt.Println("All Done")
}
