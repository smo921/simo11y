package generator

import (
	"ar/internal/generator/rand"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
)

type Metric struct {
	Name, Value, _type string
	Tags               []string
}

type metricDefinition struct {
	name, _type string
}

type metricFactory struct {
	metrics []metricDefinition
}

func MetricStream(done chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)

		statsd, err := statsd.New("127.0.0.1:12345")
		if err != nil {
			log.Fatal(err)
		}

		metrics := newMetricFactory(5, statsd)
		for {
			select {
			case <-done:
				return
			default:
			}
			metrics.SendRandomMetric(statsd)
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return out
}

func NewMetric(input string) Metric {
	m := Metric{}
	parts := strings.Split(input, "|")
	metricInfo := strings.Split(parts[0], ":")
	m.Name = metricInfo[0]
	m.Value = metricInfo[1]
	m._type = parts[1]
	m.Tags = strings.Split(parts[2], ",")
	return m
}

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

// SendRandomMetric to statsd
func (mf *metricFactory) SendRandomMetric(stats *statsd.Client) error {
	var err error
	floatVal := func() float64 { return float64(rand.SeededRand.Int()%10 + 1) }

	def := mf.metrics[rand.SeededRand.Int()%len(mf.metrics)]
	tags := randomTags()

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

func randomType() string {
	var types = []string{"c", "d", "g", "h", "s", "ms"}
	return types[rand.SeededRand.Int()%len(types)]
}

func (m Metric) ToByteSlice() []byte {
	str := fmt.Sprintf("%s:%s|%s|%s\n", m.Name, m.Value, m._type, strings.Join(m.Tags, ","))
	return []byte(str)
}
