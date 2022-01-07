package filters

import "simo11y/internal/types"

// Take the first 'limit' values from 'in' channel and write to returned 'out' channel
func Take(done chan string, limit int, in <-chan types.StructuredMessage) <-chan types.StructuredMessage {
	out := make(chan types.StructuredMessage)
	go func() {
		defer close(out)
		var count int
		for {
			select {
			case <-done:
				return
			case v := <-in:
				out <- v
				count++
				if count == limit {
					return
				}
			}
		}
	}()
	return out
}
