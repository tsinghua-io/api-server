package learn

import (
	cic "github.com/tsinghua-io/api-server/adapter/tsinghua.edu.cn/cic/learn"
	old "github.com/tsinghua-io/api-server/adapter/tsinghua.edu.cn/learn"
	"github.com/tsinghua-io/api-server/model"
	"github.com/tsinghua-io/api-server/util"
	"net/http"
	"strings"
)

type Adapter struct {
	*cic.Adapter
	Old *old.Adapter
}

func New(userId, password string) (ada *Adapter, status int, errMsg error) {
	ada = new(Adapter)

	sg := util.NewStatusGroup()
	sg.Go(func(status *int, err *error) {
		ada.Adapter, *status, *err = cic.New(userId, password)
	})
	sg.Go(func(status *int, err *error) {
		ada.Old, *status, *err = old.New(userId, password)
	})
	status, errMsg = sg.Wait()

	return
}

func (ada *Adapter) Announcements(courseId string) (announcements []*model.Announcement, status int, errMsg error) {
	if strings.Contains(courseId, "-") {
		return ada.Adapter.Announcements(courseId)
	} else {
		return ada.Old.Announcements(courseId)
	}
}

func (ada *Adapter) Assignments(courseId string) (assignments []*model.Assignment, status int, errMsg error) {
	if strings.Contains(courseId, "-") {
		return ada.Adapter.Assignments(courseId)
	} else {
		return ada.Old.Assignments(courseId)
	}
}

type fatName struct {
	name string
	seq  string
	sem  string
}

func newFatName(course *model.Course) fatName {
	return fatName{course.Name, course.CourseSequence, course.Semester}
}

func (ada *Adapter) nameIdMap() (m map[fatName]string, status int, errMsg error) {
	var courses []*model.Course
	if courses, status, errMsg = ada.Old.AllAttendedList(); errMsg == nil {
		m = make(map[fatName]string)
		for _, course := range courses {
			if strings.Contains(course.Id, "-") {
				m[newFatName(course)] = course.Id
			}
		}
	}
	return
}

func replaceCourseIds(courses []*model.Course, m map[fatName]string) {
	for _, course := range courses {
		if oldId, ok := m[newFatName(course)]; ok {
			course.Id = oldId
		}
	}
}

func (ada *Adapter) Attended(semesterID string, english bool) (courses []*model.Course, status int, errMsg error) {
	var m map[fatName]string

	sg := util.NewStatusGroup()
	sg.Go(func(status *int, err *error) {
		m, *status, *err = ada.nameIdMap()
	})
	sg.Go(func(status *int, err *error) {
		courses, *status, *err = ada.Adapter.Attended(semesterID, english)
	})
	if status, errMsg = sg.Wait(); errMsg == nil {
		replaceCourseIds(courses, m)
	}

	return
}

func (ada *Adapter) NowAttended(english bool) (thisCourses []*model.Course, nextCourses []*model.Course, status int, errMsg error) {
	var m map[fatName]string

	sg := util.NewStatusGroup()
	sg.Go(func(status *int, err *error) {
		m, *status, *err = ada.nameIdMap()
	})
	sg.Go(func(status *int, err *error) {
		thisCourses, nextCourses, *status, *err = ada.Adapter.NowAttended(english)
	})
	if status, errMsg = sg.Wait(); errMsg == nil {
		replaceCourseIds(thisCourses, m)
		replaceCourseIds(nextCourses, m)
	}

	return
}

func (ada *Adapter) PastAttended(english bool) (courses []*model.Course, status int, errMsg error) {
	var m map[fatName]string

	sg := util.NewStatusGroup()
	sg.Go(func(status *int, err *error) {
		m, *status, *err = ada.nameIdMap()
	})
	sg.Go(func(status *int, err *error) {
		courses, *status, *err = ada.Adapter.PastAttended(english)
	})
	if status, errMsg = sg.Wait(); errMsg == nil {
		replaceCourseIds(courses, m)
	}

	return
}

func (ada *Adapter) AllAttended(english bool) (courses []*model.Course, status int, errMsg error) {
	var m map[fatName]string

	sg := util.NewStatusGroup()
	sg.Go(func(status *int, err *error) {
		m, *status, *err = ada.nameIdMap()
	})
	sg.Go(func(status *int, err *error) {
		courses, *status, *err = ada.Adapter.AllAttended(english)
	})
	if status, errMsg = sg.Wait(); errMsg == nil {
		replaceCourseIds(courses, m)
	}

	return
}

func (ada *Adapter) Files(courseId string) (files []*model.File, status int, errMsg error) {
	if strings.Contains(courseId, "-") {
		return ada.Adapter.Files(courseId)
	} else {
		return ada.Old.Files(courseId)
	}
}

func (ada *Adapter) Profile() (profile *model.User, status int, errMsg error) {
	return ada.Adapter.Profile()
}

func HandlerFunc(f func(http.ResponseWriter, *http.Request, *Adapter)) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		userId, password, _ := req.BasicAuth()
		if ada, status, err := New(userId, password); err != nil {
			util.Error(rw, err.Error(), status)
		} else {
			f(rw, req, ada)
		}
	})
}
