package types

import (
	"fmt"
	"strings"
)

// Metric representation
type Metric struct {
	Name, Value, Type string
	Tags              []string
}

// FromStatsd creates a Metric from a statsd message
func MetricFromStatsd(input string) Metric {
	parts := strings.Split(input, "|")
	metricInfo := strings.Split(parts[0], ":")
	return Metric{
		Name:  metricInfo[0],
		Value: metricInfo[1],
		Type:  parts[1],
		Tags:  strings.Split(parts[2], ","),
	}
}

func (m Metric) ToByteSlice() []byte {
	str := fmt.Sprintf("%s:%s|%s|%s\n", m.Name, m.Value, m.Type, strings.Join(m.Tags, ","))
	return []byte(str)
}
