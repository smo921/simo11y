package transformers

import (
	"ar/internal/types"
	"encoding/json"
)

func StructuredMessage(done chan string, in <-chan string) <-chan types.StructuredMessage {
	out := make(chan types.StructuredMessage)

	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case msg := <-in:
				if ok, json := toJSON(msg); ok {
					out <- json
				} else {
					m := make(types.StructuredMessage)
					m["log"] = msg
					out <- m
				}
			}
		}
	}()
	return out
}

// is this JSON or unstructured?
func toJSON(msg string) (bool, types.StructuredMessage) {
	var jsonMsg map[string]interface{}
	err := json.Unmarshal([]byte(msg), &jsonMsg)
	if err != nil {
		return false, nil
	}
	return true, jsonMsg
}
