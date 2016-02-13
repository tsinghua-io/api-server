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
