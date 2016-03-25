package cic

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/tsinghua-io/api-server/resource"
	"net/http"
	"strconv"
)

func AnnouncementsURL(courseId string) string {
	return fmt.Sprintf("%s/b/myCourse/notice/listForStudent/%s?pageSize=1000", BaseURL, courseId)
}

func (ada *Adapter) Announcements(courseId string, _ map[string]string, announcements *[]*resource.Announcement) (status int) {
	if announcements == nil {
		glog.Errorf("nil received")
		return http.StatusInternalServerError
	}

	url := AnnouncementsURL(courseId)
	var v struct {
		PaginationList struct {
			RecordList []struct {
				CourseNotice struct {
					Id          int64
					Title       string
					Owner       string
					RegDate     string
					CourseId    string
					MsgPriority int // 0 for normal, 1 for important.
					Detail      string
				}
			}
		}
	}

	if err := ada.GetJSON("GET", url, &v); err != nil {
		return http.StatusBadGateway
	}

	// TODO: Iterate a pointer slice?
	for _, result := range v.PaginationList.RecordList {
		announcement := &resource.Announcement{
			Id:        strconv.FormatInt(result.CourseNotice.Id, 10),
			CourseId:  result.CourseNotice.CourseId,
			Owner:     &resource.User{Name: result.CourseNotice.Owner},
			CreatedAt: result.CourseNotice.RegDate,
			Priority:  result.CourseNotice.MsgPriority,
			Title:     result.CourseNotice.Title,
			Body:      result.CourseNotice.Detail,
		}
		*announcements = append(*announcements, announcement)
	}

	return http.StatusOK
}
