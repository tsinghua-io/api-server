package util

import (
	"encoding/json"
	"net/http"
	"sync"
)

func NotFound(rw http.ResponseWriter, _ *http.Request) {
	Error(rw, "Not Found", http.StatusNotFound)
}

func Error(rw http.ResponseWriter, err string, code int) {
	v := map[string]string{"message": err}
	rw.WriteHeader(code)
	json.NewEncoder(rw).Encode(v)
}

type StatusGroup struct {
	*sync.WaitGroup
	Status int
	sync.Mutex
}

func NewStatusGroup() *StatusGroup {
	return &StatusGroup{Status: http.StatusOK}
}

func (sg *StatusGroup) Done(status int) {
	sg.Lock()
	if sg.Status == http.StatusOK && status != http.StatusOK {
		if status == 0 {
			// Usally an early exit.
			sg.Status = http.StatusInternalServerError
		} else {
			sg.Status = status
		}
	}
	sg.Unlock()
	sg.WaitGroup.Done()
}

func (sg *StatusGroup) Wait() int {
	sg.WaitGroup.Wait()
	return sg.Status
}
