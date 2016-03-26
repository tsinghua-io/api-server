package learn

import (
	"github.com/tsinghua-io/api-server/adapter/cic/learn"
	"github.com/tsinghua-io/api-server/adapter/learn"
	"github.com/tsinghua-io/api-server/resource"
	"net/http"
	"strings"
)

type MixedAdapter struct {
	old *old.OldAdapter
	cic *cic.CicAdapter
}

func Login(username string, password string) (ada *MixedAdapter, status int) {
	var oldAda *old.OldAdapter
	oldStatus := make(chan int, 1)
	go func() {
		ada, status := old.Login(username, password)
		oldAda = ada
		oldStatus <- status
	}()

	cicAda, cicStatus := cic.Login(username, password)
	if cicStatus != http.StatusOK {
		return nil, cicStatus
	}
	if status := <-oldStatus; status != http.StatusOK {
		return nil, status
	}

	// We are safe.
	return &MixedAdapter{old: oldAda, cic: cicAda}, http.StatusOK
}

func (ada *MixedAdapter) Profile(username string, params map[string]string) (user *resource.User, status int) {
	return ada.cic.Profile(username, params)
}

func (ada *MixedAdapter) Attended(username string, params map[string]string) (courses []*resource.Course, status int) {
	cic2old := make(chan map[string]string, 1)

	go func() {
		courses, status := ada.old.Attended(username, params)
		if status != http.StatusOK {
			cic2old <- make(map[string]string)
		}
		cic2old <- old.CourseIdMap(courses)
	}()

	if courses, status = ada.cic.Attended(username, params); status != http.StatusOK {
		return nil, status
	}

	// We are safe, substitute ids.
	mapping := <-cic2old
	for _, course := range courses {
		if oldId := mapping[course.Id]; oldId != "" {
			course.Id = oldId
		}
	}
	return courses, http.StatusOK
}

func (ada *MixedAdapter) CourseAnnouncements(courseId string, params map[string]string) (announcements []*resource.Announcement, status int) {
	if strings.Contains(courseId, "-") {
		return ada.cic.CourseAnnouncements(courseId, params)
	} else {
		return ada.old.CourseAnnouncements(courseId, params)
	}
}

func (ada *MixedAdapter) CourseFiles(courseId string, params map[string]string) (files []*resource.File, status int) {
	if strings.Contains(courseId, "-") {
		return ada.cic.CourseFiles(courseId, params)
	} else {
		return ada.old.CourseFiles(courseId, params)
	}
}

func (ada *MixedAdapter) CourseHomework(courseId string, params map[string]string) (files []*resource.Homework, status int) {
	if strings.Contains(courseId, "-") {
		return ada.cic.CourseHomework(courseId, params)
	} else {
		return ada.old.CourseHomework(courseId, params)
	}
}
