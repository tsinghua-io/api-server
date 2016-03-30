package util

import (
	"fmt"
	"net/http"
	"sync"
)

type StatusGroup struct {
	sync.WaitGroup
	sync.Mutex
	Status int
	Err    error
}

func NewStatusGroup() *StatusGroup {
	return &StatusGroup{Status: http.StatusOK}
}

func (sg *StatusGroup) Done(status int, err error) {
	if status == 0 {
		// An early exit.
		status = http.StatusInternalServerError
		err = fmt.Errorf("Unkown errors occur in an goroutine.")
	}

	sg.Lock()
	if sg.Err == nil && err != nil {
		sg.Status = status
		sg.Err = err
	}
	sg.Unlock()
	sg.WaitGroup.Done()
}

func (sg *StatusGroup) Wait() (int, error) {
	sg.WaitGroup.Wait()
	return sg.Status, sg.Err
}
