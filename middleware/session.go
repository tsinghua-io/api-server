package middleware

import (
	"github.com/gorilla/context"
	"net/http"
)

type UserSession struct {
	LoginName  string
	LoginPass  string
	Session    string
	CicSession string
}

// GetUserSession get user session infomation from redis.
func GetUserSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	context.Set(r, "userSession", UserSession{
		LoginName: "nxf12",
		LoginPass: "hihihi",
	})
}
