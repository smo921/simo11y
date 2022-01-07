package generator

import (
	"fmt"
	"log"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"

	"simo11y/internal/generator/rand"
)

type metricDefinition struct {
	name, _type string
}

type metricFactory struct {
	metrics []metricDefinition
}

// Generate a stream of random metrics
func MetricStream(done chan string, Mutater string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)

		statsd, err := statsd.New(Mutater)
		if err != nil {
			log.Fatal(err)
		}

		metrics := newMetricFactory(5, statsd)
		for {
			select {
			case <-done:
				return
			case <-time.After(100 * time.Millisecond):
				metrics.SendRandomMetric(statsd)
			}
		}
	}()
	return out
}

// newMetricFactory of num random metrics.  Guarantees that a specific number of umique metric
// names/types will be created and returned by the metricFactory
func newMetricFactory(num int, client *statsd.Client) *metricFactory {
	mf := &metricFactory{
		metrics: make([]metricDefinition, num),
	}
	for x := range mf.metrics {
		mf.metrics[x] = metricDefinition{
			name:  rand.MetricName(),
			_type: randomType(),
		}
	}
	return mf
}

// SendRandomMetric from the metric factory to statsd
func (mf *metricFactory) SendRandomMetric(stats *statsd.Client) error {
	var err error
	floatVal := func() float64 { return float64(rand.SeededRand.Int()%10 + 1) }

	def := mf.metrics[rand.SeededRand.Int()%len(mf.metrics)]
	tags := randomTags()

	// call the correct method based on the type of metric
	switch def._type {
	case "c": // count
		err = stats.Incr(def.name, tags, floatVal())
	case "d": // distribution
		err = stats.Distribution(def.name, floatVal(), tags, floatVal())
	case "g": // gauge
		err = stats.Gauge(def.name, floatVal(), tags, floatVal())
	case "h": // histogram
		err = stats.Histogram(def.name, floatVal(), tags, floatVal())
	case "s": // set
		err = stats.Set(def.name, rand.String(10, rand.Charset), tags, floatVal())
	case "ms": // timing
		err = stats.TimeInMilliseconds(def.name, floatVal(), tags, floatVal())
	default:
		err = fmt.Errorf("Unknown metric type: '%s'", def._type)
	}
	return err
}

// generate a random metric "type"
func randomType() string {
	var types = []string{"c", "d", "g", "h", "s", "ms"}
	return types[rand.SeededRand.Int()%len(types)]
}
