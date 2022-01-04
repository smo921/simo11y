package consumers

import "fmt"
import "ar/internal/types"

// Basic consumer of messages from a channel
func Message(done <-chan string, in <-chan types.Message) <-chan types.Message {
	out := make(chan types.Message)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case msg, open := <-in:
				if !open {
					return
				}
				fmt.Println("\nConsumed message:", msg)
			}
		}
	}()
	return out
}

// Basic consumer that uses structured messages
func StructuredMessage(done <-chan string, in <-chan types.StructuredMessage) <-chan types.StructuredMessage {
	out := make(chan types.StructuredMessage)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case msg, open := <-in:
				if !open {
					return
				}
				fmt.Printf("\nConsumed message: %v\n", msg)
			}
		}
	}()
	return out
}
