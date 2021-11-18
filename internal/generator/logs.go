package generator

import "fmt"
import "ar/internal/generator/rand"

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
	accounts := NewAccountLogger(3)
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
