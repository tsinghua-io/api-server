// Resources Agent

package agent

import (
	"encoding/json"
	"github.com/golang/glog"
	"github.com/gorilla/context"
	"github.com/tsinghua-io/api-server/adapter/old"
	//"github.com/tsinghua-io/api-server/adapter/cic"
	"github.com/gorilla/mux"
	"github.com/tsinghua-io/api-server/webapp"
	"net/http"
	"reflect"
)

type handlerSpec struct {
	methodName string
	httpMethod string
	getArgs    func(http.ResponseWriter, *http.Request) []interface{}
}

type userAgent struct {
	UrlMap map[string]handlerSpec
}

var UserAgent = userAgent{
	UrlMap: map[string]handlerSpec{
		"/users/me":                         handlerSpec{"PersonalInfo", "GET", getArgsMe},
		"/users/me/attending":               handlerSpec{"Attending", "GET", getArgsMe},
		"/users/me/attended":                handlerSpec{"Attended", "GET", getArgsMe},
		"/courses/{courseId}/announcements": handlerSpec{"Announcements", "GET", getArgsCourse},
		"/courses/{courseId}/files":         handlerSpec{"Files", "GET", getArgsCourse},
		"/courses/{courseId}/homeworks":     handlerSpec{"Homeworks", "GET", getArgsCourse},
	},
}

func getArgsMe(w http.ResponseWriter, r *http.Request) []interface{} {
	return []interface{}{}
}

func getArgsCourse(w http.ResponseWriter, r *http.Request) []interface{} {
	vars := mux.Vars(r)
	courseId := vars["courseId"]
	return []interface{}{courseId}
}

func (useragent *userAgent) BindRoute(app *webapp.WebApp) {
	for path, handlerSpec := range useragent.UrlMap {
		app.HandleFunc(path, useragent.GenerateHandler(handlerSpec.methodName,
			handlerSpec.getArgs)).Methods(handlerSpec.httpMethod)
	}
}

func (useragent *userAgent) GenerateHandler(methodName string,
	getArgs func(http.ResponseWriter, *http.Request) []interface{}) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, ok := context.GetOk(r, "session")
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ada := old.New(session.([]*http.Cookie))

		// get the arguments
		args := getArgs(w, r)
		values := make([]reflect.Value, len(args))
		for i, arg := range args {
			values[i] = reflect.ValueOf(arg)
		}

		// call the actual handler in the adapter
		returnVals := reflect.ValueOf(ada).MethodByName(methodName).Call(values)

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
