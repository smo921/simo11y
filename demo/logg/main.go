package main

import "fmt"

import (
	"ar/internal/consumers"
	logGenerator "ar/internal/generator/logs"
	"ar/internal/generator/rand"
	"ar/internal/outputs"
	"ar/internal/sources"
	"ar/internal/transformers"
	"ar/internal/types"
)

const broker = "localhost:9092"
const topic = "demo_topic"
const keyToMutate = "foo"

func main() {
	fmt.Println("Starting")
	done := make(chan string)
	defer close(done)

	go func() {
		<-consumers.Structured(done,
			search(done, keyToMutate, rand.SeededRand.Int()%10,
				sources.Kafka(done, broker, topic, "search_demo"),
			),
		)
	}()

	// BLOCKING: Generate random log messages and write them to kafka
	outputs.Kafka(done, broker, topic,
		transformers.Mutate(done, MutateRandomField,
			transformers.LogHash(done, "logHash",
				transformers.StructuredMessage(done,
					logGenerator.SteadyStream(done, 2, logGenerator.Messages(done)),
				),
			),
		),
	)
}

func MutateRandomField(m types.StructuredMessage) types.StructuredMessage {
	m[keyToMutate] = rand.SeededRand.Int() % 10
	return m
}

func search(done chan string, key string, value int, in <-chan types.StructuredMessage) <-chan types.StructuredMessage {
	out := make(chan types.StructuredMessage)
	fmt.Printf("Search started for %s==%d\n", key, value)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case msg := <-in:
				if _, ok := msg[key]; ok {
					if v, ok := msg[key].(float64); ok {
						v := int(v)
						if value == v {
							out <- msg
						}
					}
				}
			}
		}
	}()
	return out
}
