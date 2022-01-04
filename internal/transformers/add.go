package transformers

import "ar/internal/types"

type Callback func(types.StructuredMessage) types.StructuredMessage

func Add(done chan string, cb Callback, in <-chan types.StructuredMessage) <-chan types.StructuredMessage {
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
