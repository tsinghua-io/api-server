package util

import (
	"encoding/json"
	"github.com/golang/glog"
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
	} else {
		rw.WriteHeader(status)
		if err := json.NewEncoder(rw).Encode(v); err != nil {
			// Too late, just log it.
			glog.Errorf("JSON encoding error: %s", err)
		}
	}
}

func HeadersHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		rw.Header().Add("Vary", "Authorization")
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
