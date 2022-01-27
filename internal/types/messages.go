package types

import (
	"encoding/json"
	"fmt"
	"strings"
)

type CompressedMessages []byte
type Message string
type StructuredMessage map[string]interface{}
type StructuredMessages []StructuredMessage

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
		// TODO: return error
		fmt.Printf("Error converting structured message to byte array: %s", err)
	}
	return rawLog
}

func (m StructuredMessage) Size() int {
	return len(m.Raw())
}

func (m StructuredMessages) Raw() []byte {
	rawLog, err := json.Marshal(m)
	if err != nil {
		// TODO: return error
		fmt.Printf("Error converting structured message to byte array: %s", err)
	}
	return rawLog
}

func ReadStructuredMessages(data []byte) (*StructuredMessages, error) {
	msgs := &StructuredMessages{}
	err := json.Unmarshal(data, msgs)
	return msgs, err
}
