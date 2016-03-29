package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/gorilla/handlers"
	"github.com/tsinghua-io/api-server/api"
	"github.com/tsinghua-io/api-server/resource"
	"github.com/tsinghua-io/api-server/util"
	"net/http"
	"strconv"
)

func main() {
	host := flag.String("host", "", "Host of the server")
	port := flag.Int("port", 8000, "Port of the server")
	flag.Parse()

	api := api.New(
		handlers.CompressHandler,
		util.ContentTypeHandler,
		util.NewLimiter(60, 10).Handler(),
	)

	api.AddResource("/users/me", resource.Profile)
	api.AddResource("/users/me/attended", resource.Attended)

	addr := *host + ":" + strconv.Itoa(*port)
	glog.Infof("Starting server on %s", addr)
	err := http.ListenAndServe(addr, api)
	glog.Fatalf("Shutting down: %s", err)
}
