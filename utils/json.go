package utils

import "encoding/json"

func IsJSON(buf []byte) bool {
	var js map[string]interface{}
	return json.Unmarshal(buf, &js) == nil
}
