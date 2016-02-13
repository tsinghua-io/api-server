package agent

import (
	"encoding/json"
	"fmt"
	"crypto/sha256"
	"github.com/golang/glog"
	"github.com/gorilla/context"
	"github.com/tsinghua-io/api-server/adapter/old"
	"github.com/tsinghua-io/api-server/adapter/cic"
	"gopkg.in/redis.v3"
	"net/http"
)

const (
	OldSessionTimeout = 0
	CicSessionTimeout = 0
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

		oldKey := userKey + ":old"
		cicKey := userKey + ":cic"

		oldSession, err := client.Get(oldKey).Result()
		if err == redis.Nil {
			oldSession = ""
		} else if err != nil {
			glog.Warningln("Error when fetching session from redis: \n", err)
			oldSession = ""
		}

		if oldSession == "" {
			cookies, err := old.Login(loginName, loginPass)
			if err != nil {
				if err.Error() == "Bad credentials." {
					w.WriteHeader(http.StatusUnauthorized)
				} else {
					w.WriteHeader(http.StatusBadGateway)
				}
				return
			}
			// save the cookies
			j, _ := json.Marshal(cookies)
			err = client.Set(oldKey, string(j), OldSessionTimeout).Err()
			if err != nil {
				glog.Warningln("Error when setting session to redis: \n", err)
			}
			context.Set(r, "oldSession", cookies)
		} else {
			var cookies []*http.Cookie
			if err := json.Unmarshal([]byte(oldSession), &cookies); err != nil {
				panic(err)
			}
			context.Set(r, "oldSession", cookies)
		}

		cicSession, err := client.Get(cicKey).Result()
		if err == redis.Nil {
			cicSession = ""
		} else if err != nil {
			glog.Warningln("Error when fetching session from redis: \n", err)
			cicSession = ""
		}


		if cicSession == "" {
			cookies, err := cic.Login(loginName, loginPass)
			if err != nil {
				if err.Error() == "Bad credentials." {
					w.WriteHeader(http.StatusUnauthorized)
				} else {
					w.WriteHeader(http.StatusBadGateway)
				}
				return
			}
			// save the cookies
			j, _ := json.Marshal(cookies)
			err = client.Set(cicKey, string(j), CicSessionTimeout).Err()
			if err != nil {
				glog.Warningln("Error when setting session to redis: \n", err)
			}
			context.Set(r, "cicSession", cookies)
		} else {
			var cookies []*http.Cookie
			if err := json.Unmarshal([]byte(oldSession), &cookies); err != nil {
				panic(err)
			}
			context.Set(r, "cicSession", cookies)
		}

		context.Set(r, "clearSession", func(cic bool) bool {
			var err error
			if cic {
				err = client.Set(cicKey, "", 0).Err()
			} else {
				err = client.Set(oldKey, "", 0).Err()
			}

			if err != nil {
				glog.Warningln("Error when setting session to redis: \n", err)
				return false
			}
			return true
		})
		// Call the original handler
		h.ServeHTTP(w, r)
	})
}
