package main

import "fmt"

import "ar/internal/generator"

func main() {
	fmt.Println("Starting")
	done := make(chan string)
	source := generator.LogStream(done, 5, generator.LogMessages(done))
	ch1, ch2 := tee(done, source)
	consumer1(done, ch1)
	consumer2(done, ch2)
	<-done
	fmt.Println("All Done")
}

func consumer1(done chan string, in <-chan string) {
	// consume until last message is read
	go func() {
		defer close(done)
		for {
			msg, open := <-in
			if !open {
				break
			}
			fmt.Println("\nConsumor 1 message:", msg)
		}
	}()
}

func consumer2(done chan string, in <-chan string) {
	// consume until last message is read
	go func() {
		defer close(done)
		for {
			msg, open := <-in
			if !open {
				break
			}
			fmt.Println("\nConsumor 2 message:", msg)
		}
	}()
}

func tee(done <-chan string, in <-chan string) (<-chan string, <-chan string) {
	ch1 := make(chan string)
	ch2 := make(chan string)
	go func() {
		defer close(ch1)
		defer close(ch2)
		for {
			select {
			case <-done:
				return
			case message := <-in:
				ch1 <- message
				ch2 <- message
			}
		}
	}()

	return ch1, ch2
}
