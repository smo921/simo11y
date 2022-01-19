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
		"http://localhost:6060", // local godoc web server - `godoc -http=:6060`
		"https://www.google.com/",
		"https://en.wikipedia.org/wiki/Main_Page",
		"https://www.reddit.com/",
		"https://duckduckgo.com/",
	}
	for {
		//for i := range urls {
		i := 0
		url := urls[i]
		span := tracer.StartSpan(fmt.Sprintf("get.%s", url))

		// Create a child of it, computing the time needed to fetch a url.
		child := tracer.StartSpan(url, tracer.ChildOf(span.Context()))
		_, err := http.Get(url)
		child.SetTag(ext.ResourceName, url)

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
