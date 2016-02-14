// Adapters

package adapter

import (
	"github.com/tsinghua-io/api-server/resource"
)

type Adapter interface {
	PersonalInfo() (user *resource.User, status int)

	Attending() (courses []*resource.Course, status int)
	Attended() (courses []*resource.Course, status int)

	Announcements(courseId string) (announcements []*resource.Announcement, status int)
	Files(courseId string) (files []*resource.File, status int)
	Homeworks(courseId string) (homeworks []*resource.Homework, status int)
}
