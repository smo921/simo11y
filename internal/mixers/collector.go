package mixers

import (
	"simo11y/internal/types"
)

// Collector takes messages from the in channel and returns messages grouped into slices
func Collector(done chan string, msgLimit, maxSize int, in <-chan types.StructuredMessage) <-chan types.StructuredMessages {
	out := make(chan types.StructuredMessages)
	go func() {
		defer close(out)
		var totalSize int
		collection := make(types.StructuredMessages, 0, msgLimit)
		for {
			select {
			case <-done:
				return
			case msg := <-in:
				// add message to collection and output when full
				collection = append(collection, msg)
				totalSize += msg.Size()
				if (maxSize > 0 && totalSize > maxSize) || (len(collection) > msgLimit) {
					out <- collection
					// Reset collection after writing to out channel
					collection = make(types.StructuredMessages, 0, msgLimit)
					totalSize = 0
				}
			}
		}
	}()
	return out
}
