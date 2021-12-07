package main

import "fmt"

import "ar/internal/consumers"
import "ar/internal/generator/logs"

const numMessages = 20

func main() {
	fmt.Println("Starting")
	done := make(chan string)
	source := logs.SlowStream(done, numMessages, logs.LogMessages(done))
	ch1, ch2 := tee(done, source)
	consumers.Basic(done, ch1)
	consumers.Basic(done, ch2)
	<-done
	fmt.Println("All Done")
}

func tee(done <-chan string, in <-chan string) (_, _ <-chan string) {
	ch1 := make(chan string)
	ch2 := make(chan string)
	go func() {
		defer close(ch1)
		defer close(ch2)
		for {
			select {
			case <-done:
				return
			case message := <-in:
				ch1 <- message
				ch2 <- message
			}
		}
	}()

	return ch1, ch2
}
