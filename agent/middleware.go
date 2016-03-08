package agent

import (
	"github.com/gorilla/context"
	"github.com/tsinghua-io/api-server/adapter/mixed"
	"net/http"
)

// SetContentType middleware set the Content-Type header of the response
func SetContentType(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		h.ServeHTTP(w, r)
	})
}

// GetUserSession middleware login in the user to the learning web.
// The username and password should be in the HTTP basic auth header of the request.
func GetUserSession(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loginName, loginPass, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ada, status := mixed.Login(loginName, loginPass)
		if status != http.StatusOK {
			w.WriteHeader(status)
			return
		}

		context.Set(r, "adapter", ada)

		// Call the original handler
		h.ServeHTTP(w, r)
	})
}
