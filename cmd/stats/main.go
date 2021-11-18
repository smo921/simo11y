package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"

	"ar/internal/generator"
	"ar/internal/generator/rand"
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
			statsd.Incr(metrics.RandomMetric(), randomTags(), float64(rand.SeededRand.Int()%10+1))
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
					metricStream <- newMetric(line)
				}
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
				_, err := out.Write(m)
				if err != nil {
					fmt.Println("Error", err)
					return
				}
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

func randomTags() []string {
	t := []string{"env:local"}
	t = append(t, fmt.Sprintf("account:%d", rand.SeededRand.Int()%10000))
	t = append(t, "instance:"+rand.String(10, rand.CharsetLower))
	return t
}
