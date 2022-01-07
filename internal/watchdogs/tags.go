package watchdogs

import (
	"simo11y/internal/types"
	"fmt"
	"strings"
)

type Taggregator struct {
	// map["metric_name"+"tag_name"]["tag_value"]
	tagMap map[string]map[string]int
	limit  int
}

func NewTaggregator(l int) *Taggregator {
	t := Taggregator{limit: l}
	t.tagMap = make(map[string]map[string]int)
	return &t
}

// Tagging watchdog returns true if a boundless tag is detected
func (tagger *Taggregator) Tagging(m types.Metric) types.Metric {
	for _, tag := range m.Tags {
		t := strings.Split(tag, ":")
		if len(t) != 2 { // only count name:value combinations
			continue
		}

		name := m.Name + ":" + t[0]
		value := t[1]

		if _, ok := tagger.tagMap[name]; !ok {
			tagger.tagMap[name] = make(map[string]int)
			tagger.tagMap[name][value] = 1
		} else if _, ok := tagger.tagMap[name][value]; !ok {
			tagger.tagMap[name][value] = 1
		} else if _, ok := tagger.tagMap[name][value]; ok {
			tagger.tagMap[name][value]++
		}

		if len(tagger.tagMap[name]) > tagger.limit {
			fmt.Printf("%s exceeds configured limit: %d > %d\n",
				name, len(tagger.tagMap[name]), tagger.limit)
		}
	}
	return m
}
