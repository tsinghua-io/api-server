package cic

import (
	"testing"
)

const (
	Username = ""
	Password = ""
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
	user, err := adapter.PersonalInfo()
	if err != nil {
		t.Errorf("Unable to get personal data: ", err)
		return
	}

	t.Log("Personal info received: ", user)
}
