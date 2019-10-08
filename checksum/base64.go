package checksum

import (
	"encoding/base64"
)

// Base64 simple convert to base64
func Base64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
