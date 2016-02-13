package old

import (
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
		t.Errorf("Incorrect data: %s", user)
		return
	}
}
