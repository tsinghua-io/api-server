package learn

import (
	"github.com/tsinghua-io/api-server/util"
	"net/http"
	"testing"
)

func TestSemesters(t *testing.T) {
	thisSem, nextSem, status, err := ada.Semesters()
	if err != nil {
		t.Fatalf("Failed to get semesters: %s", err)
	}

	util.ExpectStatus(t, status, http.StatusOK)

	util.ExpectDeepEqual(t, thisSem, "2015-2016-2")
	util.ExpectDeepEqual(t, nextSem, "2015-2016-3")
}
