package agent

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/gorilla/context"
	"github.com/tsinghua-io/api-server/adapter/old"
	"gopkg.in/redis.v3"
	"net/http"
)

const (
	SessionTimeout = 0
)

func GetMD5Tag(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		md5 := r.URL.Query()["md5"]
		if md5 == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		context.Set(r, "contentMD5", md5[0])
		h.ServeHTTP(w, r)
	})
}

// GetUserSession get user session infomation from redis and do login if needed.
func GetUserSession(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		loginName, loginPass, ok := r.BasicAuth()
		//fmt.Printf("login:%s, %s\n", loginName, loginPass)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Fixme: add salt to avoid look-up table attacks
		loginPassHash := fmt.Sprintf("%x", sha256.Sum256([]byte(loginPass)))
		userKey := loginName + ":" + loginPassHash

		// Fetch session cookie from redis
		client := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0, // use default DB
		})

		var clearSessionFunc = func() bool {
			var err error
			err = client.Set(userKey, "", 0).Err()

			if err != nil {
				glog.Warningln("Error when clear session to redis: \n", err)
				return false
			}
			return true
		}

		sessionJson, err := client.Get(userKey).Result()
		if err == redis.Nil {
			sessionJson = ""
		} else if err != nil {
			glog.Warningln("Error when fetching session from redis: \n", err)
			sessionJson = ""
		}

		var session []*http.Cookie
		if sessionJson == "" {
			session, err = old.Login(loginName, loginPass)
			if err != nil {
				if err.Error() == "Bad credentials." {
					w.WriteHeader(http.StatusUnauthorized)
				} else {
					w.WriteHeader(http.StatusBadGateway)
				}
				return
			}
			if j, err := json.Marshal(session); err != nil {
				glog.Warningf("Failed marshaling session: %s : %s \n", session, err)
			} else {
				if err := client.Set(userKey, j, SessionTimeout).Err(); err != nil {
					glog.Warningln("Error when setting session to redis: \n", err)
				}
			}
		} else {
			if err := json.Unmarshal([]byte(sessionJson), &session); err != nil {
				// TODO: Session cached in redis error, clear the cache
				glog.Warningf("Failed to unmarshal session: %s : %s\n", sessionJson, err)
				clearSessionFunc()
			}
		}
		context.Set(r, "session", session)

		context.Set(r, "clearSession", clearSessionFunc)

		// Call the original handler
		h.ServeHTTP(w, r)
	})
}
