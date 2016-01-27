/*
Api server in go
*/

package main

import (
	"github.com/tsinghua-io/api-server/agent"
	"github.com/tsinghua-io/api-server/middleware"
	"github.com/tsinghua-io/api-server/webapp"
	"log"
	"net/http"
)

const (
	ADDRESS = "127.0.0.1:8080"
)

func BindRoute(app *webapp.WebApp) {
	app.PreRequest(middleware.GetUserSession)
	app.PreRequest(middleware.GetMD5Tag)
	app.UseAgent(agent.UserAgent)
}

func main() {
	app := webapp.NewWebApp()
	BindRoute(app)
	http.Handle("/", app)
	err := http.ListenAndServe(ADDRESS, nil)
	if err != nil {
		log.Fatal("Error occured when lauching server: \n", err)
	}
}
