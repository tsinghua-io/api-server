// Resources Agent

package agent

import (
	"encoding/json"
	"github.com/golang/glog"
	"github.com/gorilla/context"
	"github.com/tsinghua-io/api-server/adapter/old"
	//"github.com/tsinghua-io/api-server/adapter/cic"
	"github.com/tsinghua-io/api-server/webapp"
	"net/http"
	"reflect"
)

type userAgent struct {
	URLMap map[string]string
}

var UserAgent = userAgent{
	URLMap: map[string]string{
		"/users/me":           "PersonalInfo",
		"/users/me/attending": "Attending",
		"/users/me/attended":  "Attended",
	},
}

func (useragent userAgent) BindRoute(app *webapp.WebApp) {
	for path, methodName := range useragent.URLMap {
		app.HandleFunc(path, useragent.GenerateHandler(methodName))
	}
}

func (useragent *userAgent) GenerateHandler(methodName string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, ok := context.GetOk(r, "session")
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ada := old.New(session.([]*http.Cookie))

		returnVals := reflect.ValueOf(ada).MethodByName(methodName).Call([]reflect.Value{})
		res := returnVals[0].Interface()
		status := returnVals[1].Interface().(int)
		if status != http.StatusOK {
			// clear the cookie
			clearSession, ok := context.GetOk(r, "clearSession")
			clearSessionFunc := clearSession.(func() bool)
			if !ok {
				glog.Warningln("No clearSession func in the request context.")
			}
			clearSessionFunc() // clear session of learning web
		}

		w.WriteHeader(status)
		if res != nil {
			j, _ := json.Marshal(res)
			w.Write(j)
		}

	})
}
