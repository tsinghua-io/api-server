package cic

import (
	"encoding/json"
	"fmt"
	"github.com/tsinghua-io/api-server/resource"
	"io"
)

type announcementsParser struct {
	PaginationList struct {
		RecordList []struct {
			// Status       string
			CourseNotice struct {
				Id          int64
				Title       string
				Owner       string
				RegDate     string
				CourseId    string
				MsgPriority int
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

	for _, result := range p.PaginationList.RecordList {
		announcement := &resource.Announcement{
			Id:        string(result.CourseNotice.Id),
			CourseId:  result.CourseNotice.CourseId,
			Owner:     resource.User{Name: result.CourseNotice.Owner},
			CreatedAt: result.CourseNotice.RegDate,
			Priority:  result.CourseNotice.MsgPriority,
			Title:     result.CourseNotice.Title,
			Body:      result.CourseNotice.Detail,
		}
		*announcements = append(*announcements, announcement)
	}

	return nil
}

func (adapter *CicAdapter) Announcements(course_id string) (courses []*resource.Announcement, status int) {
	return
}
