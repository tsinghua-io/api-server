package api

import (
	"github.com/gorilla/mux"
	"github.com/tsinghua-io/api-server/api/resource"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

type API struct {
	*mux.Router
	middlewares []Middleware
}

func New(middlewares ...Middleware) *API {
	return &API{
		Router:      mux.NewRouter(),
		middlewares: middlewares,
	}
}

func (api *API) Use(middlewares ...Middleware) {
	api.middlewares = append(api.middlewares, middlewares...)
}

func (api *API) AddResource(path string, r interface{}) *mux.Route {
	return api.HandleFunc(path, resource.HandlerFunc(r))
}

func (api *API) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var handler http.Handler
	handler = api.Router

	for i := len(api.middlewares) - 1; i >= 0; i-- {
		handler = api.middlewares[i](handler)
	}

	handler.ServeHTTP(rw, req)
}
