package cic

import (
	"flag"
	"github.com/golang/glog"
	"os"
	"testing"
)

const (
	username = "lisihan13"
	password = "1L2S3H@th"
)

var (
	adapter *CicAdapter
)

func TestMain(m *testing.M) {
	// flag.Set("alsologtostderr", "true")
	// flag.Set("v", "3")
	flag.Parse()

	// Login.
	cookies, err := Login(username, password)
	if err != nil {
		glog.Errorf("Failed to login to %s: %s", username, err)
		os.Exit(1)
	}
	adapter = New(cookies, "zh-CN")

	os.Exit(m.Run())
}
