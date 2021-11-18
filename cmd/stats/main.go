package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"

	"ar/internal/generator"
)

const maxBufferSize = 1024

func main() {
	done := make(chan string)
	forwarder(done, reader(done))

	go func() {
		defer close(done)

		statsd, err := statsd.New("127.0.0.1:12345")
		if err != nil {
			log.Fatal(err)
		}

		metrics := generator.NewMetricFactory(5, statsd)
		for {
			select {
			case <-done:
				return
			default:
			}
			metrics.SendRandomMetric(statsd)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	select {
	case <-done:
		fmt.Println("All Done.")
	}
}

func openPort(addr string) (*net.UDPAddr, error) {
	return net.ResolveUDPAddr("udp", addr)
}

func dial(addr string) (*net.UDPConn, error) {
	port, err := openPort(addr)
	if err != nil {
		return nil, err
	}
	return net.DialUDP("udp", nil, port)
}

func listen(addr string) (*net.UDPConn, error) {
	port, err := openPort(addr)
	if err != nil {
		return nil, err
	}
	return net.ListenUDP("udp", port)
}

func reader(done chan string) <-chan generator.Metric {
	metricStream := make(chan generator.Metric)
	go func() {
		defer close(metricStream)

		in, err := listen("127.0.0.1:12345")
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
					metricStream <- generator.NewMetric(line)
				}
			}
		}

	}()
	return metricStream
}

func forwarder(done chan string, metricStream <-chan generator.Metric) {
	out, err := dial("127.0.0.1:23456")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer out.Close()

		for {
			select {
			case <-done:
				return
			default:
				m := (<-metricStream).ToByteSlice()
				_, err := out.Write(m)
				if err != nil {
					fmt.Println("Error", err)
					return
				}
			}
		}
	}()
}
