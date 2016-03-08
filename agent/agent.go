// Resources Agent

package agent

import (
	"encoding/json"
	"github.com/golang/glog"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/tsinghua-io/api-server/webapp"
	"net/http"
	"reflect"
	"strings"
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
		"/users/me":                         handlerSpec{"Profile", "GET", getArgsMe},
		"/users/me/attended":                handlerSpec{"Attended", "GET", getArgsMe},
		"/courses/{courseId}/announcements": handlerSpec{"CourseAnnouncements", "GET", getArgsCourse},
		"/courses/{courseId}/files":         handlerSpec{"CourseFiles", "GET", getArgsCourse},
		"/courses/{courseId}/homeworks":     handlerSpec{"CourseHomework", "GET", getArgsCourse},
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
var getArgsMe = argsFromUrl("userId")

func (useragent *userAgent) BindRoute(app *webapp.WebApp) {
	for path, handlerSpec := range useragent.UrlMap {
		app.HandleFunc(path, useragent.GenerateHandler(handlerSpec.methodName,
			handlerSpec.getArgs)).Methods(handlerSpec.httpMethod)
	}
}

func (useragent *userAgent) GenerateHandler(methodName string,
	getArgs func(http.ResponseWriter, *http.Request) []interface{}) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ada, ok := context.GetOk(r, "adapter")
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// get the arguments
		args := getArgs(w, r)
		values := make([]reflect.Value, len(args)+1)
		for i, arg := range args {
			values[i] = reflect.ValueOf(arg)
		}

		// Pass url query as the final parameter
		var params = make(map[string]string)
		for k, v := range r.URL.Query() {
			params[k] = strings.Join(v, ", ")
		}
		values[len(values)-1] = reflect.ValueOf(params)

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
