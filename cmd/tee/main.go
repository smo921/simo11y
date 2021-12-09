package main

import "fmt"

import "ar/internal/consumers"
import "ar/internal/filters"
import "ar/internal/generator/logs"
import "ar/internal/transformers"
import "ar/internal/types"

const numMessages = 3

func main() {
	fmt.Println("Starting")
	done := make(chan string)
	defer close(done)
	source := transformers.StructuredMessage(done, logs.SteadyStream(done, 1, logs.LogMessages(done)))
	source = filters.Take(done, numMessages, source)
	ch1, ch2 := tee(done, source)
	c1 := consumers.Structured(done, ch1)
	c2 := consumers.Structured(done, ch2)
	<-c1
	<-c2
	fmt.Println("All Done")
}

func tee(done chan string, in <-chan types.StructuredMessage) (_, _ <-chan types.StructuredMessage) {
	ch1 := make(chan types.StructuredMessage)
	ch2 := make(chan types.StructuredMessage)
	go func() {
		defer close(ch1)
		defer close(ch2)
		for {
			select {
			case <-done:
				return
			case message, ok := <-in:
				if !ok {
					return
				}
				var ch1, ch2 = ch1, ch2
				// need to ensure all messages are sent before returning?
				for i := 0; i < 2; i++ {
					select {
					case <-done:
						return
					case ch1 <- message:
						ch1 = nil
					case ch2 <- message:
						ch2 = nil
					}
				}
			}
		}
	}()
	return ch1, ch2
}
