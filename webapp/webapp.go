package webapp

import (
	"github.com/gorilla/mux"
	"net/http"
)

type WebApp struct {
	*mux.Router
	preRequestHandlers []http.HandlerFunc
}

type routeAgent interface {
	BindRoute(app *WebApp)
}

func NewWebApp() *WebApp {
	app := &WebApp{
		Router: mux.NewRouter()}
	return app
}

// PreRequest of WebApp adds a pre-request handler.
// The lastest added middleware is called first.
func (app *WebApp) PreRequest(f http.HandlerFunc) {
	app.preRequestHandlers = append([]http.HandlerFunc{f},
		app.preRequestHandlers...)
}

func (app *WebApp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, h := range app.preRequestHandlers {
		h(w, r)
	}
	app.Router.ServeHTTP(w, r)
}

// UseAgent of WebApp adds the routes of a routeAgent to the app receiver by calling agent.BindRoute.
func (app *WebApp) UseAgent(agent routeAgent) {
	agent.BindRoute(app)
}
