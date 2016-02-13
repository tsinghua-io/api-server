/*
Api server in go
*/

package main

import (
	"github.com/golang/glog"
	"github.com/tsinghua-io/api-server/agent"
	"github.com/tsinghua-io/api-server/webapp"
	"net/http"
	"github.com/NYTimes/gziphandler"
)

const (
	ADDRESS = "127.0.0.1:8080"
)

func BindRoute(app *webapp.WebApp) {
	app.UseAgent(agent.UserAgent)
	app.UseMiddleware(agent.GetUserSession)
	//app.UseMiddleware(agent.GetMD5Tag)
}

func main() {
	app := webapp.NewWebApp()
	BindRoute(app)
	http.Handle("/", gziphandler.GzipHandler(app))
	err := http.ListenAndServe(ADDRESS, nil)
	if err != nil {
		glog.Fatalln("Error occured when lauching server: \n", err)
	}
}
