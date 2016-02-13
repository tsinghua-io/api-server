// Adapters

package adapter

import (
	"github.com/tsinghua-io/api-server/resource"
)

type Adapter interface {
	PersonalInfo() (user *resource.User, status int)

	Attending() (courses []*resource.Course, status int)
	Attended() (courses []*resource.Course, status int)

	Announcements(course_id string) (courses []*resource.Announcement, status int)
	Files(course_id string) (courses []*resource.File, status int)
	Homeworks(course_id string) (courses []*resource.Homework, status int)
}
