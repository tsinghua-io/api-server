package cic

import (
	"github.com/tsinghua-io/api-server/resource"
	"net/http"
	"testing"
)

func TestSelfProfile(t *testing.T) {
	user, status := adapter.Profile("")
	if status != http.StatusOK {
		t.Errorf("Unable to get self profile: %d", status)
		return
	}

	// Check fetched data.
	expected := resource.User{
		Id:         "2013011187",
		Name:       "李思涵",
		Department: "电子系",
		Class:      "无36 ",
		Gender:     "男",
		Email:      "lisihan969@gmail.com",
		Phone:      "18800183697",
	}
	if *user != expected {
		t.Errorf("Incorrect data: %s", user)
		return
	}
}

func BenchmarkPersonalInfo(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user, status := adapter.Profile("")
		_ = user
		_ = status
	}
}

func TestProfile(t *testing.T) {
	user, status := adapter.Profile(username)

	if user != nil {
		t.Errorf("Should return a nil User pointer.")
	} else if status != http.StatusBadRequest {
		t.Errorf("Status should be 400 BadRequest, got %d", status)
	}
}
