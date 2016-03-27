package api

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"net/http"
)

type NotFoundHandler struct{}

func (h *NotFoundHandler) ServeHTTP(rw http.ResponseWriter, _ *http.Request) {
	rw.WriteHeader(http.StatusNotFound)
}

type API struct {
	router *mux.Router
	chain  alice.Chain
}

func New(constructors ...alice.Constructor) *API {
	r := mux.NewRouter()
	r.NotFoundHandler = new(NotFoundHandler)

	return &API{
		router: r,
		chain:  alice.New(constructors...),
	}
}

func (api *API) Use(constructors ...alice.Constructor) {
	api.chain = api.chain.Append(constructors...)
}

func (api *API) AddResource(path string, r interface{}) *mux.Route {
	return api.router.HandleFunc(path, HandlerFunc(r))
}

func (api *API) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	api.chain.Then(api.router).ServeHTTP(rw, req)
}
