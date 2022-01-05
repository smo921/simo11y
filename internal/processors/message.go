package processors

import (
	"ar/internal/types"
)

// Basic consumer of messages from a channel
func Message(done <-chan string, cb types.Callback, in <-chan types.Message) <-chan types.Message {
	out := make(chan types.Message)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case msg, ok := <-in:
				if !ok {
					return
				}
				out <- cb(msg)
			}
		}
	}()
	return out
}

// Basic consumer that uses structured messages
func StructuredMessage(done chan string, cb types.StructuredCallback, in <-chan types.StructuredMessage) <-chan types.StructuredMessage {
	out := make(chan types.StructuredMessage)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case msg, ok := <-in:
				if !ok {
					return
				}
				out <- cb(msg)
			}
		}
	}()
	return out
}
