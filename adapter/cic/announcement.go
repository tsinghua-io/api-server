package cic

import (
	"encoding/json"
	"fmt"
	"github.com/tsinghua-io/api-server/resource"
	"io"
	"strconv"
	"strings"
)

// TODO: Rename to Notice?

const (
	AnnouncementsURL = BaseURL + "/b/myCourse/notice/listForStudent/{course_id}?pageSize=1000"
)

type announcementsParser struct {
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

func (p *announcementsParser) parse(r io.Reader, info interface{}, _ string) error {
	announcements, ok := info.(*[]*resource.Announcement)
	if !ok {
		return fmt.Errorf("The parser and the destination type do not match.")
	}

	dec := json.NewDecoder(r)
	if err := dec.Decode(p); err != nil {
		return err
	}

	// TODO: Iterate a pointer slice?
	for _, result := range p.PaginationList.RecordList {
		var status int
		if _, err := fmt.Sscan(result.Status, &status); err != nil {
			return err
		}

		announcement := &resource.Announcement{
			Id:        strconv.FormatInt(result.CourseNotice.Id, 10),
			CourseId:  result.CourseNotice.CourseId,
			Owner:     resource.User{Name: result.CourseNotice.Owner},
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

func (adapter *CicAdapter) Announcements(course_id string) (announcements *[]*resource.Announcement, status int) {
	announcements = &[]*resource.Announcement{}
	url := strings.Replace(AnnouncementsURL, "{course_id}", course_id, -1)
	status = adapter.FetchInfo(url, "GET", &announcementsParser{}, announcements)
	return announcements, status
}
