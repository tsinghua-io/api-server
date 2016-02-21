package cic

import (
	"github.com/tsinghua-io/api-server/resource"
	"net/http"
	"testing"
)

const (
	Username = "lisihan13"
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
	user, status := adapter.PersonalInfo("zh-CN")
	if status != http.StatusOK {
		t.Errorf("Unable to get personal data: %s", status)
		return
	}

	// Check fetched data.
	expectedUser := resource.User{
		Id:         "2013011187",
		Name:       "李思涵",
		Type:       "",
		Department: "电子系",
		Class:      "无36 ",
		Gender:     "男",
		Email:      "lisihan969@gmail.com",
		Phone:      "18800183697",
	}
	if *user != expectedUser {
		t.Errorf("Incorrect data: %s", user)
		return
	}
}

func TestAttended(t *testing.T) {
	cookies, err := Login(Username, Password)
	if err != nil {
		t.Error(err)
		return
	}

	adapter := New(cookies)
	courses, status := adapter.Attended("zh-CN")
	if status != http.StatusOK {
		t.Errorf("Unable to get attended courses: %s", status)
		return
	}

	for _, course := range courses {
		t.Log(*course)
	}
}
