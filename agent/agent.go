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
		"/users/me":                                                 handlerSpec{"PersonalInfo", "GET", getArgsMe},
		"/users/me/attending":                                       handlerSpec{"Attending", "GET", getArgsMe},
		"/users/me/attended":                                        handlerSpec{"Attended", "GET", getArgsMe},
		"/courses/{courseId}/announcements":                         handlerSpec{"Announcements", "GET", getArgsCourse},
		"/courses/{courseId}/files":                                 handlerSpec{"Files", "GET", getArgsCourse},
		"/courses/{courseId}/homeworks":                             handlerSpec{"Homeworks", "GET", getArgsCourse},
	},
}

func argsFromUrl(argNames ...string) func(http.ResponseWriter, *http.Request) []interface{} {
	return func(w http.ResponseWriter, r *http.Request) (args []interface{}) {
		vars := mux.Vars(r)
		for _, argName := range argNames {
			args = append(args, vars[argName])
		}
		return
	}
}

var getArgsCourse = argsFromUrl("courseId")
var getArgsMe = argsFromUrl()

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

		ada := old.New(session.([]*http.Cookie), "")

		// get the arguments
		args := getArgs(w, r)
		values := make([]reflect.Value, len(args))
		for i, arg := range args {
			values[i] = reflect.ValueOf(arg)
		}

		// call the actual handler in the adapter
		methodVal := reflect.ValueOf(ada).MethodByName(methodName)
		if !methodVal.IsValid() {
			glog.Errorf("The current in-use adapter is not a receiver of method named '%s': Please check the url mapping.", methodName)
			status := http.StatusInternalServerError
			w.WriteHeader(status)
			return
		}

		returnVals := methodVal.Call(values)
		res := returnVals[0].Interface()
		status := returnVals[1].Interface().(int)

		w.WriteHeader(status)
		if res != nil {
			j, _ := json.Marshal(res)
			w.Write(j)
		}

	})
}
