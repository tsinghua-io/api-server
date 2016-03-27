package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/tsinghua-io/api-server/api"
	"net/http"
	"strconv"
)

func main() {
	host := flag.String("host", "", "Host of the server")
	port := flag.Int("port", 8000, "Port of the server")
	flag.Parse()

	api := api.New()

	addr := *host + ":" + strconv.Itoa(*port)
	glog.Infof("Starting server on %s", addr)
	err := http.ListenAndServe(addr, api)
	glog.Fatalf("Shutting down: %s", err)
}
