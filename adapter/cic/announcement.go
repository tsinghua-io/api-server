package cic

import (
	"encoding/json"
	"fmt"
	"github.com/tsinghua-io/api-server/adapter"
	"github.com/tsinghua-io/api-server/resource"
	"io"
	"strconv"
	"strings"
)

const (
	courseAnnouncementsURL = BaseURL + "/b/myCourse/notice/listForStudent/{course_id}?pageSize=1000"
)

type announcementsParser struct {
	params map[string]string
	data   struct {
		PaginationList struct {
			RecordList []struct {
				Status       string // 0 for read, 1 for unread.
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
}

func (p *announcementsParser) Parse(r io.Reader, info interface{}) error {
	announcements, ok := info.(*[]*resource.Announcement)
	if !ok {
		return fmt.Errorf("The parser and the destination type do not match.")
	}

	dec := json.NewDecoder(r)
	if err := dec.Decode(&p.data); err != nil {
		return err
	}

	// TODO: Iterate a pointer slice?
	for _, result := range p.data.PaginationList.RecordList {
		var status int
		if _, err := fmt.Sscan(result.Status, &status); err != nil {
			return err
		}

		announcement := &resource.Announcement{
			Id:        strconv.FormatInt(result.CourseNotice.Id, 10),
			CourseId:  result.CourseNotice.CourseId,
			Owner:     &resource.User{Name: result.CourseNotice.Owner},
			CreatedAt: result.CourseNotice.RegDate,
			Priority:  result.CourseNotice.MsgPriority,
			Read:      status != 0,
			Title:     result.CourseNotice.Title,
			Body:      result.CourseNotice.Detail,
		}
		*announcements = append(*announcements, announcement)
	}

	return nil
}

func (ada *CicAdapter) CourseAnnouncements(course_id string, params map[string]string) (announcements []*resource.Announcement, status int) {
	URL := strings.Replace(courseAnnouncementsURL, "{course_id}", course_id, -1)
	parser := &announcementsParser{params: params}
	announcements = []*resource.Announcement{}

	status = adapter.FetchInfo(&ada.client, URL, "GET", parser, &announcements)
	return announcements, status
}
