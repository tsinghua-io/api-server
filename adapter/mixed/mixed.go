package mixed

import (
	"github.com/tsinghua-io/api-server/resource"
	"net/http"
)

type MixedAdapter struct {
}

func Login(username string, password string) (cookies []*http.Cookie, err error) {
}

func (adapter *MixedAdapter) PersonalInfo() (user *resource.User, status int) {
}

func (adapter *MixedAdapter) Attending() (courses []*resource.Course, status int) {
}

func (adapter *MixedAdapter) Attended() (courses []*resource.Course, status int) {
}

func (adapter *MixedAdapter) Announcements(courseId string) (courses []*resource.Announcement, status int) {
}

func (adapter *MixedAdapter) Files(courseId string) (courses []*resource.File, status int) {
}

func (adapter *MixedAdapter) Homeworks(courseId string) (courses []*resource.Homework, status int) {
}
