package main

import "fmt"

import "ar/internal/generator"

func main() {
	fmt.Println("Starting")
	done := make(chan string)
	consumer(done, generator.LogStream(done, 5, generator.LogMessages(done)))
	<-done
	fmt.Println("All Done")
}

func consumer(done chan string, in <-chan string) {
	// consume until last message is read
	go func() {
		defer close(done)
		for {
			msg, open := <-in
			if !open {
				break
			}
			fmt.Println("\nConsumed message:", msg)
		}
	}()
}
