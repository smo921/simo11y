package consumers

import "ar/internal/types"

func Processor(done chan string, cb types.Callback, in <-chan types.StructuredMessage) <-chan types.StructuredMessage {
	out := make(chan types.StructuredMessage)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case msg := <-in:
				out <- cb(msg)
			}
		}
	}()
	return out
}
