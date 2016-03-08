/*
Api server in go
*/

package main

import (
	"flag"
	"github.com/NYTimes/gziphandler"
	"github.com/golang/glog"
	"github.com/tsinghua-io/api-server/agent"
	"github.com/tsinghua-io/api-server/webapp"
	"net/http"
)

const (
	ADDRESS = "127.0.0.1:8080"
)

func BindRoute(app *webapp.WebApp) {
	app.UseAgent(&agent.UserAgent)
	app.UseMiddleware(agent.GetUserSession)
	app.UseMiddleware(agent.SetContentType)
	app.UseMiddleware(gziphandler.GzipHandler)
}

func main() {
	flag.Parse()

	app := webapp.NewWebApp()
	BindRoute(app)
	http.Handle("/", app)
	err := http.ListenAndServe(ADDRESS, nil)
	if err != nil {
		glog.Fatalln("Error occured when lauching server: \n", err)
	}
}
