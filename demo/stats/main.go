package main

import (
	"fmt"
	"log"
	"net"

	"ar/internal/generator"
	"ar/internal/sources"
	"ar/internal/types"
)

const src = "127.0.0.1:12345"
const dest = "127.0.0.1:23456"

func main() {
	done := make(chan string)
	defer close(done)
	forwarder(done, dest, sources.Metrics(done, src))
	metrics := sources.Metrics(done, dest)

	generator.MetricStream(done, src)
	for metric := range metrics {
		fmt.Println(metric)
	}
}

func dial(addr string) (*net.UDPConn, error) {
	port, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	return net.DialUDP("udp", nil, port)
}

func forwarder(done chan string, addr string, metricStream <-chan types.Metric) {
	out, err := dial(addr)
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
