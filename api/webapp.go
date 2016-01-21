package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

type WebApp struct {
	Router *mux.Router
}

func NewWebApp() *WebApp {
	app := &WebApp{
		Router: mux.NewRouter()}
	return app
}

func (app *WebApp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.Router.ServeHTTP(w, r)
}
