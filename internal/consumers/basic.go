package consumers

import "fmt"
import "ar/internal/types"

// Basic consumer of messages from a channel
func Basic(done <-chan string, in <-chan types.Message) {
	go func() {
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
}

// Basic consumer that uses structured messages
func Structured(done <-chan string, in <-chan types.StructuredMessage) <-chan string {
	out := make(chan string)
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
