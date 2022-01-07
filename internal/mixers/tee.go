package mixers

import "simo11y/internal/types"

func Tee(done chan string, in <-chan types.StructuredMessage) (_, _ <-chan types.StructuredMessage) {
	ch1 := make(chan types.StructuredMessage)
	ch2 := make(chan types.StructuredMessage)
	go func() {
		defer close(ch1)
		defer close(ch2)
		for {
			select {
			case <-done:
				return
			case message, ok := <-in:
				if !ok {
					return
				}
				var ch1, ch2 = ch1, ch2
				// need to ensure all messages are sent before returning?
				for i := 0; i < 2; i++ {
					select {
					case <-done:
						return
					case ch1 <- message:
						ch1 = nil
					case ch2 <- message:
						ch2 = nil
					}
				}
			}
		}
	}()
	return ch1, ch2
}
