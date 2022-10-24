package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"simo11y/internal/filters"
	"simo11y/internal/processors"
	"simo11y/internal/sources"
	"simo11y/internal/transformers"
)

const numMessages = 20
const url = "http://127.0.0.1:8080"

type exampleData struct {
	ID      int
	Message string
}

func main() {
	fmt.Println("Starting")
	done := make(chan string)
	defer close(done)

	messages := filters.Take(done, numMessages,
		processors.StructuredMessage(done, transformers.LogHash, sources.HTTP(done)),
	)

	time.Sleep(5 * time.Second) // wait for connections to be permitted

	var data = [...]exampleData{
		{1, "hello world"},
		{2, "foo bar"},
	}
	payloadData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Unable to marshal example data")
		os.Exit(1)
	}

	go func() {
		for x := 0; x < numMessages; x++ {
			reqBody := bytes.NewBuffer(payloadData)
			resp, err := http.Post(url, "application/json", reqBody)
			if err == nil {
				resp.Body.Close()
			}
		}
	}()

	for m := range messages {
		fmt.Println(m)
	}
	fmt.Println("All Done")
}
