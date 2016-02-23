package cic

import (
	"flag"
	"github.com/golang/glog"
	"github.com/tsinghua-io/api-server/resource"
	"net/http"
	"os"
	"testing"
)

const (
	Username = "lisihan13"
	Password = "1L2S3H@th"
)

var (
	adapter *CicAdapter
)

func TestMain(m *testing.M) {
	flag.Parse()

	// Login.
	cookies, err := Login(Username, Password)
	if err != nil {
		glog.Error("Failed to login: ", err)
		os.Exit(1)
	}
	glog.Info("Cookies received: ", cookies)
	adapter = New(cookies)

	os.Exit(m.Run())
}

func TestLoginFail(t *testing.T) {
	_, err := Login("InvalidUsername", "InvalidPassword")
	if err == nil {
		t.Error("Logged in using invalid username/password.")
		return
	}

	t.Log("Error received: ", err)
}

func BenchmarkLogin(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cookies, err := Login(Username, Password)
		_ = cookies
		_ = err
	}
}

func TestPersonalInfo(t *testing.T) {
	user, status := adapter.PersonalInfo("zh-CN")
	if status != http.StatusOK {
		t.Errorf("Unable to get personal data: %d", status)
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

func BenchmarkPersonalInfo(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user, status := adapter.PersonalInfo("zh-CN")
		_ = user
		_ = status
	}
}

func TestAttended(t *testing.T) {
	courses, status := adapter.Attended("zh-CN")
	if status != http.StatusOK {
		t.Errorf("Unable to get attended courses: %d", status)
		return
	}

	_ = courses
}

func BenchmarkAttended(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		courses, status := adapter.Attended("zh-CN")
		_ = courses
		_ = status
	}
}
