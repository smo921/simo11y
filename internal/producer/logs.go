package producer

import (
	"fmt"
	"time"

	"ar/internal/generator/rand"
)

var randFn = func(limit int) int {
	if limit == 0 {
		limit = 10
	}
	return (rand.SeededRand.Int() % limit)
}

func Logs(done chan string, numMessages int, messages <-chan map[string]interface{}) <-chan string {
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
