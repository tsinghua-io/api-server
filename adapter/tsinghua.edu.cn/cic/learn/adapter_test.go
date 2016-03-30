package learn

import (
	"flag"
	"github.com/golang/glog"
	"github.com/tsinghua-io/api-server/util"
	"net/http"
	"os"
	"testing"
)

var (
	ada *Adapter
)

func TestNewFail(t *testing.T) {
	if _, status, err := New("", ""); status == http.StatusOK || err == nil {
		t.Error("Logged in using no username or password.")
	}
	if _, status, err := New("", "qwerty"); status == http.StatusOK || err == nil {
		t.Error("Logged in using no username.")
	}
	if _, status, err := New(util.UserId, ""); status == http.StatusOK || err == nil {
		t.Error("Logged in using no password.")
	}
	if _, status, err := New(util.UserId, "qwerty"); status == http.StatusOK || err == nil {
		t.Error("Logged in using invalid password.")
	}
}

func BenchmarkNew(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		New(util.UserId, util.Password)
	}
}

func TestMain(m *testing.M) {
	flag.Set("logtostderr", "true")
	flag.Set("v", "3")
	flag.Parse()

	// Login.
	if tempAda, status, err := New(util.UserId, util.Password); status != http.StatusOK || err != nil {
		glog.Errorf("Failed to login to account %s: %s: %s", util.UserId, http.StatusText(status), err)
		os.Exit(1)
	} else {
		ada = tempAda
	}

	os.Exit(m.Run())
}
