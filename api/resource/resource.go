package resource

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/golang/glog"
	"net/http"
	"strings"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
	HEAD   = "HEAD"
	PATCH  = "PATCH"
)

type GetSupported interface {
	Get(*http.Request) (interface{}, int)
}

type PostSupported interface {
	Post(*http.Request) (interface{}, int)
}

type PutSupported interface {
	Put(*http.Request) (interface{}, int)
}

type DeleteSupported interface {
	Delete(*http.Request) (interface{}, int)
}

type HeadSupported interface {
	Head(*http.Request) (interface{}, int)
}

type PatchSupported interface {
	Patch(*http.Request) (interface{}, int)
}

func SupportedMethods(r interface{}) []string {
	var list []string

	if _, ok := r.(GetSupported); ok {
		list = append(list, GET)
	}
	if _, ok := r.(PostSupported); ok {
		list = append(list, POST)
	}
	if _, ok := r.(PutSupported); ok {
		list = append(list, PUT)
	}
	if _, ok := r.(DeleteSupported); ok {
		list = append(list, DELETE)
	}
	if _, ok := r.(HeadSupported); ok {
		list = append(list, HEAD)
	}
	if _, ok := r.(PatchSupported); ok {
		list = append(list, PATCH)
	}

	return list
}

func HandlerFunc(r interface{}) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {

		var handler func(*http.Request) (interface{}, int)

		switch req.Method {
		case GET:
			if r, ok := r.(GetSupported); ok {
				handler = r.Get
			}
		case POST:
			if r, ok := r.(PostSupported); ok {
				handler = r.Post
			}
		case PUT:
			if r, ok := r.(PutSupported); ok {
				handler = r.Put
			}
		case DELETE:
			if r, ok := r.(DeleteSupported); ok {
				handler = r.Delete
			}
		case HEAD:
			if r, ok := r.(HeadSupported); ok {
				handler = r.Head
			}
		case PATCH:
			if r, ok := r.(PatchSupported); ok {
				handler = r.Patch
			}
		}

		if handler == nil {
			if supported := SupportedMethods(r); len(supported) == 0 {
				// The resource does not exist at all.
				rw.WriteHeader(http.StatusNotFound)
			} else {
				// We have other methods available, tell client.
				rw.Header().Set("Allow", strings.Join(supported, ", "))
				rw.WriteHeader(http.StatusMethodNotAllowed)
			}
			return
		}
		data, code := handler(req)

		// Return gzipped data.
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		rw.Header().Set("Content-Encoding", "gzip")

		buffer := new(bytes.Buffer)
		gz := gzip.NewWriter(buffer)
		err := json.NewEncoder(gz).Encode(data)
		gz.Close()

		if err != nil {
			glog.Errorf("Failed to marshal data into JSON: %s", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Encoded, send it.
		rw.WriteHeader(code)
		rw.Write(buffer.Bytes())
	}
}
