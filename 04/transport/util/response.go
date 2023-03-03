package util

import (
	"encoding/json"
	"net/http"
)

func WriteErrResponse(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func WriteResponse(w http.ResponseWriter, statusCode int, body any) {
	w.WriteHeader(statusCode)
	b, _ := json.Marshal(body)
	w.Write(b)
}
