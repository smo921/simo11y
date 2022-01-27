package processors

import (
	"bytes"
	"compress/gzip"
	"log"
	"simo11y/internal/types"
	"time"
)

func Compressor(done <-chan string, in <-chan types.StructuredMessages) <-chan types.CompressedMessages {
	out := make(chan types.CompressedMessages)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case msg := <-in:
				var buf bytes.Buffer
				zw := gzip.NewWriter(&buf)

				// OPTIONAL: Setting the Header fields is optional.
				zw.Name = "simo11y-compresssor"
				zw.Comment = "compressed messages inside"
				zw.ModTime = time.Now()
				// END OPTIONAL

				_, err := zw.Write([]byte(msg.Raw()))
				if err != nil {
					log.Fatal(err)
				}

				if err := zw.Close(); err != nil {
					log.Fatal(err)
				}

				/*
					ratio := float64(buf.Len()) / float64(len(msg.Raw())) * 100
					fmt.Printf("Received %d bytes, compressed to %d bytes\nCompression ratio: %f%%\n",
						len(msg.Raw()), buf.Len(), ratio,
					)
				*/

				// ship the zipped data
				out <- buf.Bytes()
			}
		}
	}()
	return out
}
