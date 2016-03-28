package resource

import (
	"github.com/tsinghua-io/api-server/util"
	"net/http"
	"sort"
	"strings"
)

type Resource map[string]http.Handler

func (r Resource) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if len(r) == 0 {
		// The resource actually does not exist.
		util.NotFound(rw, req)
	} else if h, ok := r[req.Method]; !ok {
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
			util.Error(rw, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	} else {
		// We can handle it.
		h.ServeHTTP(rw, req)
	}
}
