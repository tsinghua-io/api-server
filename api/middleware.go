package api

import (
	"net/http"
)

func ContentTypeHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		h.ServeHTTP(rw, req)
	})
}
