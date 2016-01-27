package webapp

import (
	"github.com/gorilla/mux"
	"net/http"
)

type WebApp struct {
	*mux.Router
	preRequestHandlers []Handler
}

type routeAgent interface {
	BindRoute(app *WebApp)
}

func NewWebApp() *WebApp {
	app := &WebApp{
		Router: mux.NewRouter()}
	return app
}

type Handler func(w http.ResponseWriter, r *http.Request) bool

// PreRequest of WebApp adds a pre-request handler.
// The lastest added middleware is called first.
func (app *WebApp) PreRequest(f Handler) {
	app.preRequestHandlers = append([]Handler{f},
		app.preRequestHandlers...)
}

func (app *WebApp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, h := range app.preRequestHandlers {
		status := h(w, r)
		if !status {
			return
		}
	}
	app.Router.ServeHTTP(w, r)
}

// UseAgent of WebApp adds the routes of a routeAgent to the app receiver by calling agent.BindRoute.
func (app *WebApp) UseAgent(agent routeAgent) {
	agent.BindRoute(app)
}
