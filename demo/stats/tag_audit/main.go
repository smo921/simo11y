package main

import (
	"fmt"

	"ar/internal/generator"
	"ar/internal/sources"
	"ar/internal/watchdogs"
)

func main() {
	fmt.Println("Starting Tag Audit")
	done := make(chan string)
	defer close(done)
	src := fmt.Sprintf("127.0.0.1:%d", 12345)

	r := sources.Metrics(done, src)

	tagMonitor := watchdogs.NewTaggregator(2)

	generator.MetricStream(done, src)
	for m := range r {
		fmt.Println(m)
		tagMonitor.Tagging(m)
	}
}
