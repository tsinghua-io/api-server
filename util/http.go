package util

import (
	"encoding/json"
	"net/http"
)

func NotFound(rw http.ResponseWriter, _ *http.Request) {
	Error(rw, "Resource not found.", http.StatusNotFound)
}

func Error(rw http.ResponseWriter, err string, status int) {
	v := map[string]string{"message": err}
	rw.WriteHeader(status)
	json.NewEncoder(rw).Encode(v)
}

func JSON(rw http.ResponseWriter, v interface{}, status int, err error) {
	if err != nil {
		Error(rw, err.Error(), status)
	} else if body, err := json.Marshal(v); err != nil {
		Error(rw, "JSON encoding error: "+err.Error(), http.StatusInternalServerError)
	} else {
		rw.WriteHeader(status)
		rw.Write(body)
	}
}

func ContentTypeHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		h.ServeHTTP(rw, req)
	})
}

func AuthNeededHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if _, _, ok := req.BasicAuth(); !ok {
			Error(rw, "Credentials needed.", http.StatusUnauthorized)
		} else {
			h.ServeHTTP(rw, req)
		}
	})
}
