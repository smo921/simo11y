package types

import "fmt"

// TODO: Try to replace with datadog agent Span struct
// https://github.com/DataDog/datadog-agent/blob/main/pkg/trace/pb/span.proto

// Span from the datadog tracer library
type Span struct {
	Service  string             `codec:"service"`
	Name     string             `codec:"name"`
	Resource string             `codec:"resource"`
	TraceID  uint64             `codec:"trace_id"`
	SpanID   uint64             `codec:"span_id"`
	ParentID uint64             `codec:"parent_id"`
	Start    int64              `codec:"start"`
	Duration int64              `codec:"duration"`
	Error    int32              `codec:"error"`
	Meta     map[string]string  `codec:"meta"`
	Metrics  map[string]float64 `codec:"metrics"`
	Type     string             `codec:"type"`
}

func (s Span) ToString() string {
	str := fmt.Sprintf("Name: %s\n", s.Name)
	str += fmt.Sprintf("Service: %s\n", s.Service)
	str += fmt.Sprintf("Resource: %s\n", s.Resource)

	str += fmt.Sprintf("ParentID: %d\n", s.ParentID)
	str += fmt.Sprintf("TraceID: %d\n", s.TraceID)
	str += fmt.Sprintf("SpanID: %d\n", s.SpanID)

	str += fmt.Sprintf("Start: %d\n", s.Start)
	str += fmt.Sprintf("Duration: %d\n", s.Duration)

	return (str)
}
