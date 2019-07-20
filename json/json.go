package json

import (
	"bytes"
	"encoding/json"
)

// Marshal Properly Serialize JSON
func Marshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}

	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)

	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

// MarshalPretty Properly Serialize JSON with Pretty Printing
func MarshalPretty(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}

	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")

	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

// Valid Validate JSON Bytes
func Valid(bytes []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(bytes, &js) == nil
}
