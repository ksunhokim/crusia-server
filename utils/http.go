package utils

import (
	"encoding/json"
	"net/http"
)

// AS3 DFFFSADAJ
func HttpError(w http.ResponseWriter, err error, code int) {
	obj := map[string]interface{}{
		"status": code,
		"msg":    err.Error(),
	}
	json.NewEncoder(w).Encode(obj)
}

func HttpJson(w http.ResponseWriter, i interface{}) {
	obj := map[string]interface{}{
		"status": 200,
		"data":   i,
	}
	json.NewEncoder(w).Encode(obj)
}

func HttpOk(w http.ResponseWriter) {
	obj := map[string]interface{}{
		"status": 200,
	}
	json.NewEncoder(w).Encode(obj)
}
