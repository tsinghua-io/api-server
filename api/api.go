package api

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/tsinghua-io/api-server/resource"
	"github.com/tsinghua-io/api-server/util"
	"net/http"
)

type API struct {
	router *mux.Router
	chain  alice.Chain
}

func New(constructors ...alice.Constructor) *API {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(util.NotFound)

	return &API{
		router: r,
		chain:  alice.New(constructors...),
	}
}

func (api *API) Use(constructors ...alice.Constructor) {
	api.chain = api.chain.Append(constructors...)
}

func (api *API) AddResource(path string, r resource.Resource) *mux.Route {
	return api.router.Handle(path, r)
}

func (api *API) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	api.chain.Then(api.router).ServeHTTP(rw, req)
}
