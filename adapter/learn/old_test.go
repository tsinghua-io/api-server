package old

import (
	"encoding/json"
	"flag"
	"github.com/golang/glog"
	"net/http"
	"os"
	"reflect"
	"testing"
)

const (
	username = "lisihan13"
)

var (
	ada      *OldAdapter
	password = os.Getenv("thu_pass")
)

func AssertDeepEqual(t *testing.T, actual, expected interface{}) bool {
	if !reflect.DeepEqual(actual, expected) {
		actualJson, _ := json.Marshal(actual)
		expectedJson, _ := json.Marshal(expected)
		t.Errorf("Actual: %s, Expected: %s", actualJson, expectedJson)
		return false
	}
	return true
}

func TestMain(m *testing.M) {
	// flag.Set("alsologtostderr", "true")
	// flag.Set("v", "3")
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
