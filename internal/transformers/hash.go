package transformers

import (
	"crypto/sha256"
	"fmt"
	"sort"

	"simo11y/internal/types"
)

// LogHash calculates a unique hash based on the structure of a log message
func LogHash(m types.StructuredMessage) types.StructuredMessage {
	identifierKeys := [...]string{"application"}

	// get message structure
	msgStructure := messageStructure(m)
	hash := sha256.Sum256([]byte(fmt.Sprintf("%v", msgStructure)))

	var id string
	for _, key := range identifierKeys {
		id += fmt.Sprintf(",%v", m[key])
	}
	m["logHash"] = hash

	return m
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
