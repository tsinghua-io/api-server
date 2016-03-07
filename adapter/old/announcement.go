package old

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/tsinghua-io/api-server/adapter"
	"github.com/tsinghua-io/api-server/resource"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	announcementBodyURL    = BaseURL + "/MultiLanguage/public/bbs/note_reply.jsp?bbs_type=课程公告&id={id}&course_id={course_id}"
	courseAnnouncementsURL = BaseURL + "/MultiLanguage/public/bbs/getnoteid_student.jsp?course_id={course_id}"
)

type announcementBodyParser struct {
	params map[string]string
}

func (p *announcementBodyParser) Parse(r io.Reader, info interface{}) error {
	body, ok := info.(*string)
	if !ok {
		return fmt.Errorf("The parser and the destination type do not match.")
	}

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return err
	}

	tempBody, err := doc.Find("tr[height='300'] td.tr_l2").Html()
	if err != nil {
		return err
	}
	*body = tempBody
	return nil
}

type courseAnnouncementsParser struct {
	params   map[string]string
	courseId string
}

func (p *courseAnnouncementsParser) Parse(r io.Reader, info interface{}) error {
	announcements, ok := info.(*[]*resource.Announcement)
	if !ok {
		return fmt.Errorf("The parser and the destination type do not match.")
	}

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return err
	}

	rows := doc.Find("tr.tr1, tr.tr2")
	*announcements = make([]*resource.Announcement, rows.Size())

	// Parse each row.
	rows.EachWithBreak(func(i int, s *goquery.Selection) bool {
		annc := &resource.Announcement{}
		annc.CourseId = p.courseId

		if err = parseAnnouncementsRow(s, annc); err != nil {
			return false
		}
		(*announcements)[i] = annc
		return true
	})
	return err
}

func parseAnnouncementsRow(s *goquery.Selection, annc *resource.Announcement) error {
	infos := s.Find("td").Map(func(i int, tdSelection *goquery.Selection) (info string) {
		info, _ = tdSelection.Html()
		return
	})
	if len(infos) < 5 {
		return fmt.Errorf("No enough cols: %d", len(infos))
	}

	hrefSelection := s.Find("td a")
	href, exists := hrefSelection.Attr("href")
	if !exists {
		return fmt.Errorf("href not found")
	}
	hrefUrl, err := url.Parse(href)
	if err != nil {
		return err
	}

	// id.
	id := hrefUrl.Query().Get("id")
	if id == "" {
		return fmt.Errorf("Cannot get id from %s", hrefUrl)
	}

	// title & priority
	var priority int
	var title string
	switch hrefSelection.Nodes[0].FirstChild.Type {
	case html.TextNode:
		priority = 0
		title, err = hrefSelection.Html()
		if err != nil {
			return err
		}
	default:
		priority = 1
		title, err = hrefSelection.Children().Html()
		if err != nil {
			return err
		}
	}

	// We are safe.
	annc.Id = id
	annc.Owner = &resource.User{Name: infos[2]}
	annc.CreatedAt = infos[3]
	annc.Priority = priority
	annc.Read = true
	annc.Title = title

	return nil
}

func (ada *OldAdapter) CourseAnnouncements(courseId string, params map[string]string) (announcements []*resource.Announcement, status int) {
	URL := strings.Replace(courseAnnouncementsURL, "{course_id}", courseId, -1)
	parser := &courseAnnouncementsParser{params: params, courseId: courseId}

	if status = adapter.FetchInfo(&ada.client, URL, "GET", parser, &announcements); status != http.StatusOK {
		return nil, status
	}
	count := len(announcements)
	statuses := make(chan int, count)

	URL = strings.Replace(announcementBodyURL, "{course_id}", courseId, -1)
	for _, annc := range announcements {
		annc := annc

		go func() {
			URL := strings.Replace(URL, "{id}", annc.Id, -1)
			parser := &announcementBodyParser{params: params}
			statuses <- adapter.FetchInfo(&ada.client, URL, "GET", parser, &annc.Body)
		}()
	}

	// Drain the channel.
	for i := 0; i < count; i++ {
		if status = <-statuses; status != http.StatusOK {
			return nil, status
		}
	}

	return
}
