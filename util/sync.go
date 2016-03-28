package util

import (
	"net/http"
	"sync"
)

type StatusGroup struct {
	*sync.WaitGroup
	Status int
	Err    error
	sync.Mutex
}

func NewStatusGroup() *StatusGroup {
	return &StatusGroup{Status: http.StatusOK}
}

func (sg *StatusGroup) Done(status int, err error) {
	sg.Lock()
	if sg.Err == nil && err != nil {
		if status == 0 {
			// Usually caused by an early exit.
			status = http.StatusInternalServerError
		}
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
