package main

import (
	"fmt"
	"time"

	"ar/internal/logs"
	"ar/internal/logs/rand"
)

var randFn = func(limit int) int {
	if limit == 0 {
		limit = 10
	}
	return (rand.SeededRand.Int() % limit)
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
	accounts := logs.NewAccountLogger(3)
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
