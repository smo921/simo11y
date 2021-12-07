package consumers

import "fmt"

// Basic consumer of messages from a channel
func Basic(done chan string, in <-chan string) {
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