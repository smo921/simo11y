package generator

import "github.com/DataDog/datadog-go/v5/statsd"

import "ar/internal/generator/rand"

type metricDefinition struct {
	name, _type string
	action      interface{} //func(string, []string, float64) error
}

type MetricFactory struct {
	metrics []metricDefinition
}

func NewMetricFactory(num int, client *statsd.Client) *MetricFactory {
	mf := &MetricFactory{
		metrics: make([]metricDefinition, num),
	}
	for x := range mf.metrics {
		mf.metrics[x] = metricDefinition{
			name:   rand.MetricName(),
			_type:  "c",
			action: client.Incr,
		}
	}
	return mf
}

func (mf *MetricFactory) RandomMetric() string {
	return mf.metrics[rand.SeededRand.Int()%len(mf.metrics)].name
}
