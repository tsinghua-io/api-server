package learn

import (
	"github.com/tsinghua-io/api-server/util"
	"net/http"
	"testing"
)

func TestSemesters(t *testing.T) {
	var current, next string
	if status := ada.Semesters(&current, &next); status != http.StatusOK {
		t.Fatalf("Unable to get semesters: %s", http.StatusText(status))
	}

	util.AssertDeepEqual(t, current, "2015-2016-2")
	util.AssertDeepEqual(t, next, "2015-2016-3")
}
