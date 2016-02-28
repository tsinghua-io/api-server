package cic

import (
	"testing"
)

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
		cookies, err := Login(username, password)
		_ = cookies
		_ = err
	}
}
