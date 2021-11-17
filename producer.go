package main

import (
	"fmt"
	"math/rand"
	"time"
)

var randFn = func(limit int) int {
	if limit == 0 {
		limit = 10
	}
	return (seededRand.Int() % limit)
}

func producer(done chan string, numMessages int, messages <-chan map[string]interface{}) <-chan string {
	var count int
	out := make(chan string)
	go func() {
		defer close(out)
		start := time.Now()
		sleepTime := 1
		for {
			count++

			select {
			case <-done:
				return
			case <-time.After(time.Duration(sleepTime) * time.Second):
				diff := time.Now().Sub(start).Seconds()
				out <- fmt.Sprintf("(%f sec): %v", diff, <-messages)
			}
			if count == numMessages {
				out <- fmt.Sprint("producer finished")
				return
			}
			sleepTime = randFn(10)
			if randFn(100) < -1 { // slow message rate disabled
				sleepTime += 20
			}
		}
	}()
	return out
}

func messages(done chan string) <-chan map[string]interface{} {
	accounts := newAccountLogger(3)
	fmt.Println(accounts.dump())

	out := make(chan map[string]interface{})
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			default:
			}
			out <- accounts.randomLog()
		}
	}()
	return out
}

func newMessage() map[string]interface{} {
	// create a json like message with a random number of top level and nested attributes
	topLevel := rand.Int()%10 + 3
	message := make(map[string]interface{})

	for i := 0; i < topLevel; i++ {
		key := fmt.Sprintf("topLevelAttribute_%d", i)
		entry := make(map[string]interface{})
		for j := 0; j < rand.Int()%10+3; j++ {
			key2 := fmt.Sprintf("logAttribute_%d", j)
			entry[key2] = randomLogEntry()
		}
		message[key] = entry
	}
	return message
}
