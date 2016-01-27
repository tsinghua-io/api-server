// middlewares for all requests

package middleware

import (
	"github.com/gorilla/context"
	"net/http"
	"strings"
)

type UserSession struct {
	LoginName  string
	LoginPass  string
	Session    string
	CicSession string
}

// GetUserSession get user session infomation from redis.
func GetMD5Tag(w http.ResponseWriter, r *http.Request) bool {
	md5 := r.URL.Query()["md5"]
	if md5 == nil {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}

	context.Set(r, "contentMD5", md5[0])

	return true
}

func GetUserSession(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	token := r.URL.Query()["token"]
	if token == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}

	loginNameAndPass := strings.SplitN(token[0], "@", 2)
	// TODO: fetch session cookie from redis

	context.Set(r, "userSession", UserSession{
		LoginName: loginNameAndPass[0],
		LoginPass: loginNameAndPass[1],
	})
	return true
}
