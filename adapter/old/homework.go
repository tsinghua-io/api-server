package old

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/golang/glog"
	"github.com/tsinghua-io/api-server/adapter"
	"github.com/tsinghua-io/api-server/resource"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	homeworkDetailURL = BaseURL + "/MultiLanguage/lesson/student/hom_wk_detail.jsp?id={homework_id}&course_id={course_id}"
	submissionURL     = BaseURL + "/MultiLanguage/lesson/student/hom_wk_view.jsp?id={homework_id}&course_id={course_id}"
	courseHomeworkURL = BaseURL + "/MultiLanguage/lesson/student/hom_wk_brw.jsp?course_id={course_id}&language=cn"
)

func NewFloat32(f float32) *float32 {
	return &f
}

func (ada *OldAdapter) fillAttachment(attachment *resource.Attachment) {
	if attachment == nil {
		return
	}

	filename, size, err := ParseFileInfo(ada.client, attachment.DownloadUrl)
	if err != nil {
		glog.Errorf("Failed to parse file info from %s: %s", attachment.DownloadUrl, err)
		return
	}
	attachment.Filename = filename
	attachment.Size = size

	return
}

type homeworkDetailParser struct {
	params map[string]string
}

func (p *homeworkDetailParser) Parse(r io.Reader, info interface{}) error {
	hw, ok := info.(*resource.Homework)
	if !ok {
		return fmt.Errorf("The parser and the destination type do not match.")
	}

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return err
	}

	// Body.
	bodyTr := doc.Find("table#table_box tr:nth-child(2)")
	if bodyTr.Size() == 0 {
		return fmt.Errorf("Failed to find the body of the homework.")
	}
	body := bodyTr.Find("td ~ td").First().Children().Text()

	// Attachment.
	var attachment *resource.Attachment
	hrefSelection := bodyTr.Next().Find("td a")
	if fileHref, _ := hrefSelection.Attr("href"); fileHref != "" {
		attachment = &resource.Attachment{DownloadUrl: BaseURL + fileHref}
	}

	// We are safe.
	hw.Body = body
	hw.Attachment = attachment
	return nil
}

type submissionParser struct {
	params map[string]string
}

func (p *submissionParser) Parse(r io.Reader, info interface{}) error {
	submission, ok := info.(*resource.Submission)
	if !ok {
		return fmt.Errorf("The parser and the destination type do not match.")
	}

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return err
	}

	tableSelection := doc.Find("#table_box")
	if tableSelection.Length() == 0 {
		return fmt.Errorf("#table_box not found in submission page")
	}

	// Body
	infoTr := tableSelection.Find("tr:nth-child(2)")
	submission.Body = infoTr.Find("td.title+td").Children().Text()
	// Attachment
	infoTr = infoTr.Next()
	hrefSelection := infoTr.Find("td a")
	if fileHref, _ := hrefSelection.Attr("href"); fileHref != "" {
		submission.Attachment = &resource.Attachment{DownloadUrl: BaseURL + fileHref}
	}

	// Markuser and Markedat
	infoTr = infoTr.Next().Next()
	infos := infoTr.Find("td.title+td").Map(func(i int, s *goquery.Selection) (info string) {
		info = strings.TrimSpace(s.Text())
		return
	})
	// Mark
	infoTr = infoTr.Next() // score tr
	score, _ := infoTr.Find("td.title+td").Html()

	if markedBy := infos[0]; markedBy != "" {
		submission.MarkedBy = &resource.User{
			Name: strings.TrimSpace(infos[0]),
		}
	}
	if markedAt := infos[1]; markedAt != "null" {
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
	comment := infoTr.Find("td.title+td").Children().Text()
	submission.Comment = strings.TrimSpace(comment)
	// CommentAttachment
	infoTr = infoTr.Next()
	hrefSelection = infoTr.Find("td a")
	if fileHref, _ := hrefSelection.Attr("href"); fileHref != "" {
		submission.CommentAttachment = &resource.Attachment{DownloadUrl: BaseURL + fileHref}
	}

	return nil
}

type courseHomeworkParser struct {
	params   map[string]string
	courseId string
}

func (p *courseHomeworkParser) Parse(r io.Reader, info interface{}) error {
	homework, ok := info.(*[]*resource.Homework)
	if !ok {
		return fmt.Errorf("The parser and the destination type do not match.")
	}

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return err
	}

	trs := doc.Find("tr.tr1, tr.tr2")
	*homework = make([]*resource.Homework, trs.Size())

	trs.Each(func(i int, s *goquery.Selection) {
		infos := s.Find("td").Map(func(i int, tdSelection *goquery.Selection) (info string) {
			info = strings.TrimSpace(tdSelection.Text())
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
			title := hrefSelection.Text()

			hw := &resource.Homework{
				Id:       homeworkId,
				CourseId: p.courseId,
				BeginAt:  infos[1],
				DueAt:    infos[2],
				Title:    title,
			}
			if infos[3] != "尚未提交" {
				hw.Submissions = append(hw.Submissions, &resource.Submission{HomeworkId: homeworkId})
			}
			(*homework)[i] = hw
		}
	})

	return nil
}

func (ada *OldAdapter) CourseHomework(courseId string, params map[string]string) (homework []*resource.Homework, status int) {
	URL := strings.Replace(courseHomeworkURL, "{course_id}", courseId, -1)
	homeworkDetailURL := strings.Replace(homeworkDetailURL, "{course_id}", courseId, -1)
	submissionURL := strings.Replace(submissionURL, "{course_id}", courseId, -1)
	parser := &courseHomeworkParser{params: params, courseId: courseId}

	if status = adapter.FetchInfo(&ada.client, URL, "GET", parser, &homework); status != http.StatusOK {
		return nil, status
	}
	count := len(homework)
	chanSize := 2 * count
	statuses := make(chan int, chanSize)

	for _, hw := range homework {
		hw := hw

		go func() {
			URL := strings.Replace(homeworkDetailURL, "{homework_id}", hw.Id, -1)
			parser := &homeworkDetailParser{params: params}

			status := adapter.FetchInfo(&ada.client, URL, "GET", parser, hw)
			ada.fillAttachment(hw.Attachment)
			statuses <- status
		}()

		go func() {
			if len(hw.Submissions) == 0 {
				statuses <- http.StatusOK
				return
			}
			submission := hw.Submissions[0]
			URL := strings.Replace(submissionURL, "{homework_id}", hw.Id, -1)
			parser := &submissionParser{params: params}

			status := adapter.FetchInfo(&ada.client, URL, "GET", parser, submission)
			ada.fillAttachment(submission.Attachment)
			ada.fillAttachment(submission.CommentAttachment)
			statuses <- status
		}()
	}

	// Drain the channel.
	for i := 0; i < chanSize; i++ {
		if status = <-statuses; status != http.StatusOK {
			return nil, status
		}
	}

	return
}
