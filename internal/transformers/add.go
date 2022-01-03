package transformers

import "ar/internal/types"

func Add(done chan string, dest string, value interface{}, in <-chan types.StructuredMessage) <-chan types.StructuredMessage {
	out := make(chan types.StructuredMessage)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case msg := <-in:
				msg[dest] = value
				out <- msg
			}
		}
	}()
	return out
}
