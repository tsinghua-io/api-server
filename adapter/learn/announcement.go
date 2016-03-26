package learn

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/golang/glog"
	"github.com/tsinghua-io/api-server/resource"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
)

func AnnouncementBodyURL(courseId, id string) string {
	return fmt.Sprintf("%s/MultiLanguage/public/bbs/note_reply.jsp?bbs_type=课程公告&course_id=%s&id=%s", BaseURL, courseId, id)
}

func AnnouncementsURL(courseId string) string {
	return fmt.Sprintf("%s/MultiLanguage/public/bbs/getnoteid_student.jsp?course_id=%s", BaseURL, courseId)
}

func (ada *Adapter) AnnouncementBody(courseId, id string, _ map[string]string, body *string) (status int) {
	if body == nil {
		glog.Errorf("nil received")
		return http.StatusInternalServerError
	}

	url := AnnouncementBodyURL(courseId, id)
	doc, err := ada.GetDocument(url)
	if err != nil {
		return http.StatusBadGateway
	}

	tempBody, err := doc.Find("tr[height='300'] td.tr_l2").Html()
	if err != nil {
		return http.StatusBadGateway
	}
	*body = tempBody

	return http.StatusOK
}

func (ada *Adapter) Announcements(courseId string, _ map[string]string, announcements *[]*resource.Announcement) (status int) {
	if announcements == nil {
		glog.Errorf("nil received")
		return http.StatusInternalServerError
	}

	url := AnnouncementsURL(courseId)
	doc, err := ada.GetDocument(url)
	if err != nil {
		return http.StatusBadGateway
	}

	rows := doc.Find("tr.tr1, tr.tr2")
	*announcements = make([]*resource.Announcement, rows.Size())

	// Parse each row.
	statuses := make(chan int, 1)

	rows.Each(func(i int, s *goquery.Selection) {
		annc := &resource.Announcement{CourseId: courseId}
		go func() {
			statuses <- ada.parseAnnouncementsRow(s, annc)
		}()
		(*announcements)[i] = annc
	})

	status = http.StatusOK
	for i := 0; i < rows.Size(); i++ {
		if s := <-statuses; s != http.StatusOK {
			status = s
		}
	}

	return status
}

func (ada *Adapter) parseAnnouncementsRow(s *goquery.Selection, annc *resource.Announcement) (status int) {
	infos := s.Find("td").Map(func(i int, tdSelection *goquery.Selection) (info string) {
		info, _ = tdSelection.Html()
		return info
	})
	if len(infos) < 5 {
		glog.Errorf("No enough cols: %d", len(infos))
		return http.StatusBadGateway
	}

	hrefSelection := s.Find("td a")
	href, exists := hrefSelection.Attr("href")
	if !exists {
		glog.Errorf("href not found")
		return http.StatusBadGateway
	}
	hrefURL, err := url.Parse(href)
	if err != nil {
		glog.Errorf("Failed to parse href (%s): %s", href, err)
		return http.StatusBadGateway
	}

	// Id.
	id := hrefURL.Query().Get("id")
	if id == "" {
		glog.Errorf("Cannot get id from %s", hrefURL)
		return http.StatusBadGateway
	}

	// Title & priority.
	var priority int
	switch hrefSelection.Nodes[0].FirstChild.Type {
	case html.TextNode:
		priority = 0
	default:
		priority = 1
	}

	// Body.
	var body string
	if status := ada.AnnouncementBody(annc.CourseId, id, nil, &body); status != http.StatusOK {
		return status
	}

	// We are safe.
	annc.Id = id
	annc.Owner = &resource.User{Name: infos[2]}
	annc.CreatedAt = infos[3]
	annc.Priority = priority
	annc.Title = hrefSelection.Text()
	annc.Body = body

	return http.StatusOK
}
