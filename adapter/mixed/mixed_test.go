package mixed

import (
	"encoding/json"
	"flag"
	"github.com/golang/glog"
	"net/http"
	"os"
	"testing"
)

const (
	username = "lisihan13"
)

var (
	ada      *MixedAdapter
	password = os.Getenv("thu_pass")
)

func TestMain(m *testing.M) {
	flag.Set("alsologtostderr", "true")
	flag.Set("v", "3")
	flag.Parse()

	// Login.
	var status int
	ada, status = Login(username, password)
	if status != http.StatusOK {
		glog.Errorf("Failed to login to %s: %s", username, http.StatusText(status))
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func TestSelfProfile(t *testing.T) {
	user, status := ada.Profile("", nil)
	info, _ := json.Marshal(user)
	t.Logf("user: %s, status: %s", info, http.StatusText(status))
}

func TestAttended(t *testing.T) {
	courses, status := ada.Attended("", nil)
	info, _ := json.Marshal(courses)
	t.Logf("courses: %s, status: %s", info, http.StatusText(status))
}
