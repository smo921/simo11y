package main

import "fmt"

import "ar/internal/generator"
import "ar/internal/producer"

func main() {
	fmt.Println("Starting")
	done := make(chan string)
	consumer(done, producer.Logs(done, 5, generator.LogMessages(done)))
	<-done
	fmt.Println("All Done")
}
