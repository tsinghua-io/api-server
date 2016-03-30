package learn

import (
	"github.com/tsinghua-io/api-server/model"
	"github.com/tsinghua-io/api-server/util"
	"net/http"
	"testing"
)

func TestProfile(t *testing.T) {
	actual, status, err := ada.Profile()
	if err != nil {
		t.Fatalf("Failed to get Profile: %s", err)
	}

	util.ExpectStatus(t, status, http.StatusOK)

	// Check fetched data.
	expected := &model.User{
		Id:         "2013011187",
		Name:       "李思涵",
		Department: "电子系",
		Class:      "无36 ",
		Gender:     "男",
		Email:      "lisihan969@gmail.com",
		Phone:      "18800183697",
	}

	util.ExpectDeepEqual(t, actual, expected)
}

func BenchmarkPersonalInfo(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ada.Profile()
	}
}
