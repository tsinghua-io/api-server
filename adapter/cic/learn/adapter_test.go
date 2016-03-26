package learn

import (
	"flag"
	"github.com/golang/glog"
	"github.com/tsinghua-io/api-server/adapter"
	"net/http"
	"os"
	"testing"
)

var (
	ada *Adapter
)

func TestNewFail(t *testing.T) {
	if _, status := New("", ""); status == http.StatusOK {
		t.Error("Logged in using no username or password.")
	}
	if _, status := New("", "qwerty"); status == http.StatusOK {
		t.Error("Logged in using no username.")
	}
	if _, status := New("2013011187", ""); status == http.StatusOK {
		t.Error("Logged in using no password.")
	}
	if _, status := New("2013011187", "qwerty"); status == http.StatusOK {
		t.Error("Logged in using invalid password.")
	}
}

func BenchmarkNew(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		New(adapter.UserId, adapter.Password)
	}
}

func TestMain(m *testing.M) {
	flag.Set("alsologtostderr", "true")
	flag.Set("v", "3")
	flag.Parse()

	// Login.
	var status int
	ada, status = New(adapter.UserId, adapter.Password)
	if status != http.StatusOK {
		glog.Errorf("Failed to login to %s: %s", adapter.UserId, http.StatusText(status))
		os.Exit(1)
	}

	os.Exit(m.Run())
}
