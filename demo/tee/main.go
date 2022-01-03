package main

import "fmt"

import "ar/internal/consumers"
import "ar/internal/filters"
import "ar/internal/generator/logs"
import "ar/internal/mixers"
import "ar/internal/transformers"

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
