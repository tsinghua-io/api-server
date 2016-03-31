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
	host := flag.String("host", "api.tsinghua.io", "Host of the server")
	port := flag.Int("port", 443, "Port of the server")
	certFile := flag.String("cert", "", "Certificate file.")
	keyFile := flag.String("key", "", "key file.")
	windowMin := flag.Int64("window", 15, "Window size of the rate limit (in minutes).")
	rate := flag.Int("rate", 900, "Max requests per window per IP.")

	flag.Parse()

	api := api.New(
		handlers.CompressHandler,
		util.HeadersHandler,
		util.NewLimiter(*windowMin*60, *rate).Handler(),
	)

	api.AddResource("/users/me", resource.Profile)
	api.AddResource("/users/me/attended", resource.Attended)
	api.AddResource("/courses/{id}/announcements", resource.CourseAnnouncements)
	api.AddResource("/courses/{id}/files", resource.CourseFiles)
	api.AddResource("/courses/{id}/assignments", resource.CourseAssignments)

	addr := *host + ":" + strconv.Itoa(*port)
	glog.Infof("windowMin = %d, rate = %d", *windowMin, *rate)
	glog.Infof("Starting server on %s", addr)
	err := http.ListenAndServeTLS(addr, *certFile, *keyFile, api)
	glog.Fatalf("Shutting down: %s", err)
}
