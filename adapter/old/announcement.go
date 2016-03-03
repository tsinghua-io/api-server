package old

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/golang/glog"
	"github.com/tsinghua-io/api-server/resource"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
)

const (
	AnnouncementURL  = BaseURL + "/MultiLanguage/public/bbs/note_reply.jsp?bbs_type=课程公告&id={id}&course_id={course_id}"
	AnnouncementsURL = BaseURL + "/MultiLanguage/public/bbs/getnoteid_student.jsp?course_id={course_id}"
)

func (adapter *OldAdapter) announcementBody(path string) (body string) {
	path = "/MultiLanguage/public/bbs/" + path
	doc, err := adapter.getOldResponse(path, make(map[string]string))

	if err != nil {
		err = fmt.Errorf("Failed to get response from learning web: %s", err)
	} else {
		body, err = doc.Find("tr[height='300'] td.tr_l2").Html()
	}
	return
}

func (adapter *OldAdapter) parseRow(s *goquery.Selection, courseId string, annc **resource.Announcement) error {
	infos := s.Find("td").Map(func(i int, tdSelection *goquery.Selection) (info string) {
		info, _ = tdSelection.Html()
		return
	})

	hrefSelection := s.Find("td a")

	var href string
	var hrefUrl *url.URL
	var err error
	href, _ = hrefSelection.Attr("href")

	if hrefUrl, err = url.Parse(href); err != nil {
		return err
	}

	if announcementId := hrefUrl.Query().Get("id"); announcementId != "" {
		body := adapter.announcementBody(href)

		var priority int
		var title string
		switch hrefSelection.Nodes[0].FirstChild.Type {
		case html.TextNode:
			priority = 0
			title, _ = hrefSelection.Html()
		default:
			priority = 1
			title, _ = hrefSelection.Children().Html()
		}

		*annc = &resource.Announcement{
			Id:       announcementId,
			CourseId: courseId,

			Owner: &resource.User{
				Name: infos[2],
			},
			CreatedAt: infos[3],
			Priority:  priority,
			Read:      true, // We just read it.

			Title: title,
			Body:  body,
		}
		return nil
	} else {
		return fmt.Errorf("Cannot get announcement_id from %s", hrefUrl)
	}
}

func (adapter *OldAdapter) Announcements(courseId string) (announcements []*resource.Announcement, status int) {
	path := "/MultiLanguage/public/bbs/getnoteid_student.jsp?course_id=" + courseId
	doc, err := adapter.getOldResponse(path, make(map[string]string))

	if err != nil {
		glog.Errorf("Failed to get response from learning web: %s", err)
		return nil, http.StatusBadGateway
	}

	trs := doc.Find("tr.tr1, tr.tr2")
	count := trs.Size()
	announcements = make([]*resource.Announcement, count)
	errChan := make(chan error, count)

	i := 0
	trs.Each(func(_ int, s *goquery.Selection) {
		index := i
		go func() {
			errChan <- adapter.parseRow(s, courseId, &announcements[index])
		}()
		i++
	})

	for i = 0; i < count; i++ {
		if err := <-errChan; err != nil {

			return nil, http.StatusBadGateway
		}
	}

	return announcements, http.StatusOK
}
