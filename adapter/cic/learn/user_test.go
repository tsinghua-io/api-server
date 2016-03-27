package learn

import (
	"github.com/tsinghua-io/api-server/model"
	"github.com/tsinghua-io/api-server/util"
	"net/http"
	"testing"
)

func TestUser(t *testing.T) {
	var actual model.User
	if status := ada.User("", nil, &actual); status != http.StatusOK {
		t.Errorf("Unable to get self profile: %s", http.StatusText(status))
		return
	}

	// Check fetched data.
	expected := model.User{
		Id:         "2013011187",
		Name:       "李思涵",
		Department: "电子系",
		Class:      "无36 ",
		Gender:     "男",
		Email:      "lisihan969@gmail.com",
		Phone:      "18800183697",
	}

	util.AssertDeepEqual(t, actual, expected)
}

func BenchmarkPersonalInfo(b *testing.B) {
	var user model.User

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ada.User("", nil, &user)
	}
}
