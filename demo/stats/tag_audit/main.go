package main

import (
	"fmt"

	"ar/internal/generator"
	"ar/internal/sources"
	"ar/internal/watchdogs"
)

const src = "127.0.0.1:12345"

func main() {
	fmt.Println("Starting Tag Audit")
	done := make(chan string)
	defer close(done)

	r := sources.Metrics(done, src)

	tagMonitor := watchdogs.NewTaggregator(2)

	generator.MetricStream(done, src)
	for m := range r {
		fmt.Println(m)
		tagMonitor.Tagging(m)
	}
}
