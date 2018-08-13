package utils

import (
	"fmt"
	"encoding/json"
	"net/http"
)

func HttpError(w http.ResponseWriter, err error, code int) {
	fmt.Println(err)
	http.Error(w, err.Error(), code)
}

func HttpJson(w http.ResponseWriter, i interface{}) {
	err := json.NewEncoder(w).Encode(i)
	if err != nil {
		HttpError(w, err, 500)
	}
}
