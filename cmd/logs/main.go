package main

import "fmt"

import "ar/internal/consumers"
import "ar/internal/generator/logs"

const numMessages = 20

func main() {
	fmt.Println("Starting")
	done := make(chan string)
	consumers.Basic(done, logs.SteadyStream(done, numMessages, 2, logs.LogMessages(done)))
	<-done
	fmt.Println("All Done")
}
