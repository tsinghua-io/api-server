package api

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/tsinghua-io/api-server/resource"
	"net/http"
)

type API struct {
	router *mux.Router
	chain  alice.Chain
}

func New(constructors ...alice.Constructor) *API {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(resource.NotFound)

	return &API{
		router: r,
		chain:  alice.New(constructors...),
	}
}

func (api *API) Use(constructors ...alice.Constructor) {
	api.chain = api.chain.Append(constructors...)
}

func (api *API) AddResource(path string, r resource.Resource) {
	api.router.Handle(path, r)
}

func (api *API) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	api.chain.Then(api.router).ServeHTTP(rw, req)
}
