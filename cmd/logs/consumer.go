package main

import "fmt"

func consumer(in <-chan string, done chan string) {
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
