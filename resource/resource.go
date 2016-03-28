package resource

import (
	"encoding/json"
	"github.com/golang/glog"
	"net/http"
	"sort"
	"strings"
)

func NotFound(rw http.ResponseWriter, _ *http.Request) {
	Error(rw, "Not Found", http.StatusNotFound)
}

func Error(rw http.ResponseWriter, err string, code int) {
	v := map[string]string{"message": err}
	rw.WriteHeader(code)
	json.NewEncoder(rw).Encode(v)
}

type Resource map[string]func(*http.Request) (interface{}, int)

func (r Resource) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if len(r) == 0 {
		// The resource actually does not exist.
		NotFound(rw, req)
	} else if f, ok := r[req.Method]; ok {
		// We can handle it.
		data, code := f(req)

		if body, err := json.Marshal(data); err != nil {
			glog.Errorf("Failed to marshal data into JSON: %s", err)
			rw.WriteHeader(http.StatusInternalServerError)
		} else {
			rw.WriteHeader(code)
			rw.Write(body)
		}
	} else {
		// We don't support this method.
		var allow []string
		for k := range r {
			allow = append(allow, k)
		}
		sort.Strings(allow)
		rw.Header().Set("Allow", strings.Join(allow, ", "))
		if req.Method == "OPTIONS" {
			rw.WriteHeader(http.StatusOK)
		} else {
			Error(rw, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}
