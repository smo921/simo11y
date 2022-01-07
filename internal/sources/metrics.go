package sources

import (
	"fmt"
	"log"
	"net"
	"strings"

	"simo11y/internal/types"
)

const maxBufferSize = 1024

func listen(addr string) (*net.UDPConn, error) {
	port, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	return net.ListenUDP("udp", port)
}

// reads metrcis from addr and sends to returned channel
func Metrics(done chan string, addr string) <-chan types.Metric {
	metricStream := make(chan types.Metric)
	go func() {
		defer close(metricStream)

		in, err := listen(addr)
		if err != nil {
			log.Fatal(err)
		}
		defer in.Close()

		buffer := make([]byte, maxBufferSize)

		for {
			select {
			case <-done:
				return
			default:
				_, err := in.Read(buffer)
				if err != nil {
					fmt.Println("Error", err)
					return
				}

				// split on new line and write to channel
				i := strings.LastIndex(string(buffer), "\n")
				lines := strings.Split(string(buffer[:i]), "\n")

				buffer = make([]byte, maxBufferSize)
				for _, line := range lines {
					metricStream <- types.MetricFromStatsd(line)
				}
			}
		}

	}()
	return metricStream
}
