package types

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Message string
type StructuredMessage map[string]interface{}

func (m StructuredMessage) Fetch(path string) (interface{}, error) {
	steps := strings.Split(path, ".")
	numSteps := len(steps)
	location := m
	for _, step := range steps {
		if numSteps > 1 {
			v, ok := location[step].(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("unable to traverse path '%s', failed on step '%s'", path, step)
			}
			numSteps--
			location = v
		} else {
			return location[step], nil
		}
	}
	return nil, nil
}

func (m StructuredMessage) Raw() []byte {
	rawLog, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("Error converting structured message to byte array: %s", err)
	}
	return rawLog
}
