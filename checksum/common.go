package checksum

import (
	"encoding/base64"
	"encoding/json"
)

// Base64 simple convert to base64
func Base64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// PrettyPrint convert struct to pretty string
func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
