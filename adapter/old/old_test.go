package old

import (
	"encoding/json"
	"github.com/tsinghua-io/api-server/resource"
	"net/http"
	"testing"
)

const (
	Password = ""
	Username = "nxf12"
)

func TestLoginSuccuss(t *testing.T) {
	cookies, err := Login(Username, Password)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("Cookies received: ", cookies)
}

func TestLoginFail(t *testing.T) {
	_, err := Login("InvalidUsername", "InvalidPassword")
	if err == nil {
		t.Error("Logged in using invalid username/password.")
		return
	}

	t.Log("Error received: ", err)
}

func TestPersonalInfo(t *testing.T) {
	cookies, err := Login(Username, Password)
	if err != nil {
		t.Error(err)
		return
	}

	adapter := New(cookies)
	user, status := adapter.PersonalInfo()
	if status != http.StatusOK {
		t.Errorf("Unable to get personal data: %s", err)
		return
	}

	// Check fetched data.
	expectedUser := resource.User{
		Id:         "2012011067",
		Name:       "宁雪妃",
		Type:       "本科生",
		Department: "",
		Gender:     "女",
		Email:      "1175267294@qq.com",
		Phone:      "13120098897",
	}
	if *user != expectedUser {
		t.Errorf("Incorrect data: expected %s, get %s", expectedUser, user)
		return
	}
}

func testEq(a, b []string) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestAttendingIds(t *testing.T) {
	cookies, err := Login(Username, Password)
	if err != nil {
		t.Error(err)
		return
	}

	adapter := New(cookies)
	attendingIdList, err := adapter.attendingIds()
	if err != nil {
		t.Errorf("Unable to get attending course id list: %s", err)
		return
	}

	expectedIdList := []string{
		"133593", "133106", "131792", "133107", "131777",
	}
	if !testEq(attendingIdList, expectedIdList) {
		t.Errorf("Incorrect data: excpected %s, get %s", expectedIdList, attendingIdList)
		return
	}

}

func TestCourseInfo(t *testing.T) {
	cookies, err := Login(Username, Password)
	if err != nil {
		t.Error(err)
		return
	}

	adapter := New(cookies)
	course, err := adapter.courseInfo("133593")
	if err != nil {
		t.Errorf("Unable to get course info: %s", err)
		return
	}

	j, _ := json.Marshal(*course)
	t.Logf("%s", j)
	t.Errorf("") // for debuging now
}

func TestAttending(t *testing.T) {
	cookies, err := Login(Username, Password)
	if err != nil {
		t.Error(err)
		return
	}

	adapter := New(cookies)
	courseList, status := adapter.Attending()
	if status != http.StatusOK {
		t.Errorf("Unable to get attending course info: %s", err)
		return
	}

	j, _ := json.Marshal(courseList)
	t.Logf("%s", j)
	t.Errorf("") // for debugging now
}
