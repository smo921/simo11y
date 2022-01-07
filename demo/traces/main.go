package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"simo11y/internal/sources"
)

const tracePort = "localhost:8126"

func main() {
	fmt.Println("Starting trace demo")
	done := make(chan string)
	defer close(done)

	traces := sources.Traces(done, tracePort)

	tracer.Start(
		tracer.WithEnv("demo"),
		tracer.WithService("tracer-demo"),
		tracer.WithServiceVersion("abc123"),
	)
	defer tracer.Stop()

	go makeTraces()
	for trace := range traces {
		fmt.Println(trace)
	}
}

func makeTraces() {
	urls := []string{
		"https://www.google.com/",
		"https://en.wikipedia.org/wiki/Main_Page",
		"https://www.reddit.com/",
		"https://duckduckgo.com/",
	}
	for {
		fmt.Println("Making span")
		span := tracer.StartSpan("get.data")

		// Perform an operation.
		//for i := range urls {
		i := 1
		_, err := http.Get(urls[i])
		fmt.Println("http.GET", urls[i])
		// Create a child of it, computing the time needed to read a file.
		child := tracer.StartSpan(urls[i], tracer.ChildOf(span.Context()))
		child.SetTag(ext.ResourceName, urls[i])

		// We may finish the child span using the returned error. If it's
		// nil, it will be disregarded.
		child.Finish(tracer.WithError(err))
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(5 * time.Second)
		//}

		span.Finish()

		time.Sleep(5 * time.Second)
	}
}
