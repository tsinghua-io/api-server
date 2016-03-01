package old

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/golang/glog"
	"github.com/tsinghua-io/api-server/resource"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func (adapter *OldAdapter) parseHomeworkInfo(href string) (body string, attachment *resource.Attachment) {
	path := "/MultiLanguage/lesson/student/" + href
	doc, err := adapter.getOldResponse(path, make(map[string]string))

	if err != nil {
		glog.Errorf("Failed to get response from learning web: %s", err)
	} else {
		bodyTr := doc.Find("table#table_box tr:nth-child(2)")
		body, _ = bodyTr.Find("td ~ td").First().Children().Html()

		// attachment
		hrefSelection := bodyTr.Next().Find("td a")
		if fileHref, _ := hrefSelection.Attr("href"); fileHref != "" {
			attachment = adapter.parseAttachmentInfo(fileHref)
		}

	}
	return
}

func (adapter *OldAdapter) parseAttachmentInfo(fileHref string) (attachment *resource.Attachment) {
	filename, size := adapter.parseFileInfo(fileHref)
	attachment = &resource.Attachment{
		Filename:    filename,
		Size:        size,
		DownloadUrl: fileHref,
	}
	return
}

func (adapter *OldAdapter) parseSubmissions(courseId string, homeworkId string) (submissions []*resource.Submission) {
	query := url.Values{}
	query.Set("course_id", courseId)
	query.Set("id", homeworkId)
	path := "/MultiLanguage/lesson/student/hom_wk_view.jsp?" + query.Encode()
	doc, err := adapter.getOldResponse(path, make(map[string]string))

	if err != nil {
		glog.Errorf("Failed to get response from learning web: %s", err)
	} else {
		tableSelection := doc.Find("#table_box")
		if tableSelection.Length() == 0 {
			return
		}
		var submission = &resource.Submission{}
		// Body
		infoTr := tableSelection.Find("tr:nth-child(2)")
		submission.Body, _ = infoTr.Find("td.title+td").Children().Html()
		// Attachment
		infoTr = infoTr.Next()
		hrefSelection := infoTr.Find("td a")
		if fileHref, _ := hrefSelection.Attr("href"); fileHref != "" {
			submission.Attachment = adapter.parseAttachmentInfo(fileHref)
		}

		// Markuser and Markedat
		infoTr = infoTr.Next().Next()
		infos := infoTr.Find("td.title+td").Map(func(i int, s *goquery.Selection) (info string) {
			info, _ = s.Html()
			return
		})
		// Mark
		infoTr = infoTr.Next() // score tr
		score, _ := infoTr.Find("td.title+td").Html()

		submission.MarkedBy = &resource.User{
			Name: strings.TrimSpace(infos[0]),
		}
		if markedAt := strings.TrimSpace(infos[1]); markedAt != "null" {
			submission.MarkedAt = markedAt
			if score := strings.TrimSpace(score); strings.Contains(score, "分") {
				if mark, err := strconv.ParseFloat(strings.TrimSuffix(score, "分"), 32); err == nil {
					f := float32(mark)
					submission.Mark = &f
				}
			}
		}

		// Comment
		infoTr = infoTr.Next()
		comment, _ := infoTr.Find("td.title+td").Children().Html()
		submission.Comment = strings.TrimSpace(comment)
		// CommentAttachment
		infoTr = infoTr.Next()
		hrefSelection = infoTr.Find("td a")
		if fileHref, _ := hrefSelection.Attr("href"); fileHref != "" {
			submission.CommentAttachment = adapter.parseAttachmentInfo(fileHref)
		}
		submissions = append(submissions, submission)
	}
	return
}

func (adapter *OldAdapter) Homeworks(courseId string) (homeworks []*resource.Homework, status int) {
	path := "/MultiLanguage/lesson/student/hom_wk_brw.jsp?course_id=" + courseId
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

			if homeworkId := hrefUrl.Query().Get("id"); homeworkId != "" {
				title, _ := hrefSelection.Html()

				homework := &resource.Homework{
					Id:        homeworkId,
					CourseId:  courseId,
					Title:     title,
					CreatedAt: infos[1],
					BeginAt:   infos[1],
					DueAt:     infos[2],
				}
				homework.Body, homework.Attachment = adapter.parseHomeworkInfo(href)
				homework.Submissions = adapter.parseSubmissions(courseId, homeworkId)
				homeworks = append(homeworks, homework)
			}
		})
		status = http.StatusOK
	}
	return
}
