package learn

import (
	"fmt"
	"github.com/tsinghua-io/api-server/model"
	"net/http"
	"strconv"
)

func AnnouncementsURL(courseId string) string {
	return fmt.Sprintf("%s/b/myCourse/notice/listForStudent/%s?pageSize=1000", BaseURL, courseId)
}

func (ada *Adapter) Announcements(courseId string) (announcements []*model.Announcement, status int) {
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
		return nil, http.StatusBadGateway
	}

	// TODO: Iterate a pointer slice?
	for _, result := range v.PaginationList.RecordList {
		announcement := &model.Announcement{
			Id:        strconv.FormatInt(result.CourseNotice.Id, 10),
			CourseId:  result.CourseNotice.CourseId,
			Owner:     &model.User{Name: result.CourseNotice.Owner},
			CreatedAt: result.CourseNotice.RegDate,
			Priority:  result.CourseNotice.MsgPriority,
			Title:     result.CourseNotice.Title,
			Body:      result.CourseNotice.Detail,
		}
		announcements = append(announcements, announcement)
	}

	return announcements, http.StatusOK
}
