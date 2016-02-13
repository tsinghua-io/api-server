package webapp

import (
	"github.com/gorilla/mux"
	"net/http"
)

type WebApp struct {
	*mux.Router
	handler http.Handler
}

type routeAgent interface {
	BindRoute(app *WebApp)
}

func NewWebApp() *WebApp {
	app := &WebApp{
		Router: mux.NewRouter()}
	app.handler = app.Router
	return app
}

type Middleware func(http.Handler) http.Handler

// PreRequest of WebApp adds a pre-request handler.
// The lastest added middleware is called first.
func (app *WebApp) UseMiddleware(f Middleware) {
	app.handler = f(app.handler)
}

func (app *WebApp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.handler.ServeHTTP(w, r)
}

// UseAgent of WebApp adds the routes of a routeAgent to the app receiver by calling agent.BindRoute.
func (app *WebApp) UseAgent(agent routeAgent) {
	agent.BindRoute(app)
}
