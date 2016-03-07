package old

import (
	"net/http"
	"testing"
)

func TestLoginFail(t *testing.T) {
	_, status := Login("InvalidUsername", "InvalidPassword")
	if status == http.StatusOK {
		t.Error("Logged in using invalid username/password.")
		return
	}

	t.Log("Error received: ", http.StatusText(status))
}

func BenchmarkLogin(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ada, status := Login(username, password)
		_ = ada
		_ = status
	}
}
