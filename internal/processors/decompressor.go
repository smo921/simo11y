package processors

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"simo11y/internal/types"
)

// Decompressor takes a stream of compressed messages and decompresses them onto an output channel
func Decompressor(done <-chan string, in <-chan types.CompressedMessages) <-chan types.StructuredMessages {
	out := make(chan types.StructuredMessages)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case data := <-in:
				zr, err := gzip.NewReader(bytes.NewReader(data))
				if err != nil {
					fmt.Println("Gzip new reader error:", err)
					continue
				}

				buf, err := ioutil.ReadAll(zr)
				if err != nil {
					fmt.Println("Error decompressing data:", err)
					continue
				}

				if err := zr.Close(); err != nil {
					fmt.Println("Error closing gzip reader:", err)
				}

				// ship the unzipped data
				msgs, err := types.ReadStructuredMessages(buf)
				if err != nil {
					fmt.Println("Error reading decompressed messages:", err)
					continue
				}
				out <- *msgs
			}
		}
	}()
	return out
}
