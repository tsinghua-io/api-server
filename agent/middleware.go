package agent

import (
	"github.com/gorilla/context"
	"github.com/tsinghua-io/api-server/adapter/old"
	"net/http"
)

const (
	SessionTimeout = 0
)

func SetContentType(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		h.ServeHTTP(w, r)
	})
}

// GetUserSession login in the user.
func GetUserSession(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loginName, loginPass, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		session, status := old.Login(loginName, loginPass)
		if status != http.StatusOK {
			w.WriteHeader(status)
			return
		}

		context.Set(r, "session", session)

		// Call the original handler
		h.ServeHTTP(w, r)
	})
}
