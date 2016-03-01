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

func (adapter *OldAdapter) Announcements(courseId string) (announcements []*resource.Announcement, status int) {
	path := "/MultiLanguage/public/bbs/getnoteid_student.jsp?course_id=" + courseId
	doc, err := adapter.getOldResponse(path, make(map[string]string))

	status = http.StatusBadGateway

	if err != nil {
		glog.Errorf("Failed to get response from learning web: %s", err)
	} else {
		trs := doc.Find("tr.tr1, tr.tr2")
		trs.Each(func(i int, s *goquery.Selection) {
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
				return
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

				announcements = append(announcements, &resource.Announcement{
					Id:       announcementId,
					CourseId: courseId,

					Owner: &resource.User{
						Name: infos[2],
					},
					CreatedAt: infos[3],
					Priority:  priority,

					Title: title,
					Body:  body,
				})
			}
		})
		status = http.StatusOK
	}
	return
}
