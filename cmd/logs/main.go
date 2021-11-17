package main

import "fmt"

func main() {
	fmt.Println("Starting")
	done := make(chan string)
	consumer(producer(done, 5, messages(done)), done)
	<-done
	fmt.Println("All Done")
}
