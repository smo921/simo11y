package transformers

import (
	"crypto/sha256"
	"fmt"
	"sort"

	"ar/internal/types"
)

// LogHash calculates a unique hash based on the structure of a log message
func LogHash(done chan string, dest string, in <-chan types.StructuredMessage) <-chan types.StructuredMessage {
	out := make(chan types.StructuredMessage)
	identifierKeys := [...]string{"application"}

	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case msg := <-in:
				// get message structure
				msgStructure := messageStructure(msg)
				hash := sha256.Sum256([]byte(fmt.Sprintf("%v", msgStructure)))

				var id string
				for _, key := range identifierKeys {
					id += fmt.Sprintf(",%v", msg[key])
				}
				// return id[1:], hash
				msg[dest] = hash
				out <- msg
			}
		}
	}()
	return out
}

func messageStructure(msg types.StructuredMessage) types.StructuredMessage {
	ret := make(map[string]interface{})

	keys := make([]string, 0, len(msg))
	for k := range msg {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		switch v := msg[k].(type) {
		case string:
			if ok, val := toJSON(fmt.Sprintf("%v", v)); ok {
				ret[k] = messageStructure(val)
			} else {
				ret[k] = "string"
			}
		case int:
			ret[k] = "int"
		case float32, float64:
			ret[k] = "float"
		case interface{}:
			ret[k] = messageStructure(v.(map[string]interface{}))
		default:
			ret[k] = "UNKNOWN_TYPE"
		}
	}
	return ret
}
