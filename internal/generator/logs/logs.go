package logs

import (
	"ar/internal/generator/rand"
	"fmt"
	"time"
)

const numAccounts = 5

var randFn = func(limit int) int {
	if limit == 0 {
		limit = 10
	}
	return (rand.SeededRand.Int() % limit)
}

// SteadyStream produces a steady stream of log messages at rate messages / sec
func SteadyStream(done chan string, numMessages int, rate int, messages <-chan map[string]interface{}) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		var count int
		for {
			count++
			select {
			case <-done:
				return
			case <-time.After(1 * time.Second):
				for i := 0; i < rate; i++ {
					out <- fmt.Sprintf("%v", <-messages)
				}
				if count == numMessages {
					out <- fmt.Sprint("producer finished")
					return
				}
			}
		}
	}()
	return out
}

// SlowStream produces a slow stream of numMessages log messages
func SlowStream(done chan string, numMessages int, messages <-chan map[string]interface{}) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		var count int
		sleepTime := 1
		start := time.Now()
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

func newLog() map[string]interface{} {
	// create a json like message with a random number of top level and nested attributes
	topLevel := rand.SeededRand.Int()%10 + 3
	message := make(map[string]interface{})

	for i := 0; i < topLevel; i++ {
		key := fmt.Sprintf("topLevelAttribute_%d", i)
		entry := make(map[string]interface{})
		for j := 0; j < rand.SeededRand.Int()%10+3; j++ {
			key2 := fmt.Sprintf("logAttribute_%d", j)
			entry[key2] = rand.LogEntry()
		}
		message[key] = entry
	}
	return message
}

func LogMessages(done chan string) <-chan map[string]interface{} {
	accounts := newAccountLogger(numAccounts)
	fmt.Println(accounts.Dump())

	out := make(chan map[string]interface{})
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			default:
			}
			out <- accounts.RandomLog()
		}
	}()
	return out
}
