package util

import (
	"encoding/json"
	"net/http"
)

func NotFound(rw http.ResponseWriter, _ *http.Request) {
	Error(rw, "Not Found", http.StatusNotFound)
}

func Error(rw http.ResponseWriter, err string, code int) {
	v := map[string]string{"message": err}
	rw.WriteHeader(code)
	json.NewEncoder(rw).Encode(v)
}
