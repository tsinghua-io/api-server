package learn

import (
	"github.com/tsinghua-io/api-server/adapter"
	"net/http"
	"testing"
)

func TestSemesters(t *testing.T) {
	var current, next string
	if status := ada.Semesters(&current, &next); status != http.StatusOK {
		t.Errorf("Unable to get semesters: %s", http.StatusText(status))
		return
	}

	adapter.AssertDeepEqual(t, current, "2015-2016-2")
	adapter.AssertDeepEqual(t, next, "2015-2016-3")
}
