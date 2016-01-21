/*
Api server in go
*/

package main

import (
	"github.com/tsinghua-io/api-server/api"
	"log"
	"net/http"
)

const (
	ADDRESS = "127.0.0.1:8080"
)

func main() {
	app := api.NewWebApp()
	app.BindRoute()
	http.Handle("/", app)
	err := http.ListenAndServe(ADDRESS, nil)
	if err != nil {
		log.Fatal("Error occured when lauching server: \n", err)
	}
}
