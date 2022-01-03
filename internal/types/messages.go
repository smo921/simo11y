package types

import (
	"encoding/json"
	"fmt"
)

type Message string
type StructuredMessage map[string]interface{}

func (m StructuredMessage) Raw() []byte {
	rawLog, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("Error converting structured message to byte array: %s", err)
	}
	return rawLog
}
