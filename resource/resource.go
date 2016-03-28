package resource

import (
	"encoding/json"
	"github.com/golang/glog"
	"github.com/gorilla/handlers"
	"net/http"
)

type Resource map[string]RESTHandlerFunc
type RESTHandlerFunc func(*http.Request) (interface{}, int)

func (r Resource) Handler() http.Handler {
	h := make(map[string]http.Handler)
	for name, f := range r {
		h[name] = RESTHandlerFunc(f)
	}
	return handlers.MethodHandler(h)
}

func (f RESTHandlerFunc) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	data, status := f(req)

	body, err := json.Marshal(data)
	if err != nil {
		glog.Errorf("Failed to marshal data into JSON: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(status)
	rw.Write(body)
}
