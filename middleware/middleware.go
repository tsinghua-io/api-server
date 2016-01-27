// middlewares for all requests

package middleware

import (
	"github.com/gorilla/context"
	"net/http"
	"strings"
	"gopkg.in/redis.v3"
	"github.com/golang/glog"
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
	loginName := loginNameAndPass[0]

	// Fetch session cookie from redis
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,  // use default DB
	})
	oldSession, err := client.Get(loginName + ":old").Result()
	if err == redis.Nil {
		oldSession = ""
	} else if err != nil {
		glog.Warningln("Error when fetching session from redis: \n", err)
		oldSession = ""
	}

	cicSession, err := client.Get(loginName + ":cic").Result()
	if err == redis.Nil {
		cicSession = ""
	} else if err != nil {
		glog.Warningln("Error when fetching session from redis: \n", err)
		cicSession = ""
	}

	context.Set(r, "setSession", func (session string, cic bool) bool {
		var key string
		if cic {
			key = loginName + ":cic"
		} else {
			key = loginName + ":old"
		}
		// fixme: should exists an expiration
		err := client.Set(key, session, 0).Err()
		if err != nil {
			glog.Warningln("Error when setting session to redis: \n", err)
			return false
		}
		return true
	})

	context.Set(r, "userSession", UserSession{
		LoginName: loginName,
		LoginPass: loginNameAndPass[1],
		Session: oldSession,
		CicSession: cicSession,
	})

	return true
}
