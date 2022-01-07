package main

import (
	"fmt"

	"simo11y/internal/consumers"
	"simo11y/internal/filters"
	"simo11y/internal/generator/logs"
	"simo11y/internal/mixers"
	"simo11y/internal/transformers"
)

const numMessages = 3

func main() {
	fmt.Println("Starting")
	done := make(chan string)
	defer close(done)
	source := transformers.StructuredMessage(done, logs.SteadyStream(done, 1, logs.Messages(done)))
	source = filters.Take(done, numMessages, source)
	ch1, ch2 := mixers.Tee(done, source)
	c1 := consumers.StructuredMessage(done, ch1)
	c2 := consumers.StructuredMessage(done, ch2)
	<-c1
	<-c2
	fmt.Println("All Done")
}
