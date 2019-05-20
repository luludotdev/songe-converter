package main

import (
	"bytes"
	"encoding/json"
)

// JSONMarshal Properly Serialize JSON
func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

// JSONMarshalPretty Properly Serialize JSON with Pretty Printing
func JSONMarshalPretty(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

// IsJSON Validate JSON Bytes
func IsJSON(bytes []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(bytes, &js) == nil
}
