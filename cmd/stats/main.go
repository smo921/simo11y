package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"

	"ar/internal/generator/rand"
)

const maxBufferSize = 1024

func main() {
	statsd, err := statsd.New("127.0.0.1:12345")
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan string)

	forwarder(done, reader(done))

	go func() {
		defer close(done)
		for {
			select {
			case <-done:
				return
			default:
			}
			statsd.Incr("foo.bar.count", []string{"env:local", "account:1234", "instance:boo-ya"}, float64(rand.SeededRand.Int()%10+1))
			time.Sleep(100 * time.Millisecond)
		}
	}()

	select {
	case <-done:
		fmt.Println("All Done.")
	}
}

type metric struct {
	name, value, tags, _type string
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

func reader(done chan string) <-chan metric {
	metricStream := make(chan metric)
	in, err := listen("127.0.0.1:12345")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer close(metricStream)
		defer in.Close()

		for {
			select {
			case <-done:
				return
			default:
				buffer := make([]byte, maxBufferSize)
				n, err := in.Read(buffer)
				if err != nil {
					fmt.Println("Error", err)
					return
				}
				fmt.Printf("Read %d bytes: %s", n, buffer)
				metricStream <- newMetric(string(buffer))
			}
		}

	}()
	return metricStream
}

func forwarder(done chan string, metricStream <-chan metric) {
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
				m := (<-metricStream).toByteSlice()
				n, err := out.Write(m)
				if err != nil {
					fmt.Println("Error", err)
					return
				}
				fmt.Printf("Wrote %d bytes\n", n)
			}
		}
	}()
}

func newMetric(input string) metric {
	m := metric{}
	parts := strings.Split(input, "|")
	metricInfo := strings.Split(parts[0], ":")
	m.name = metricInfo[0]
	m.value = metricInfo[1]
	m._type = parts[1]
	m.tags = parts[2]
	return m
}

func (m metric) toByteSlice() []byte {
	str := fmt.Sprintf("%s:%s|%s|%s\n", m.name, m.value, m._type, m.tags)
	return []byte(str)
}
