package mixers

import "simo11y/internal/types"

// TODO; handle case of closed in1/in2 channels

func Combine(done chan string, in1, in2 <-chan types.StructuredMessage) <-chan types.StructuredMessage {
	out := make(chan types.StructuredMessage)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case message := <-in1:
				out <- message
			case message := <-in2:
				out <- message
			}
		}
	}()
	return out
}
