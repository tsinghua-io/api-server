package learn

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/golang/glog"
	"github.com/tsinghua-io/api-server/resource"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func HomeworkDetailURL(courseId, homeworkId string) string {
	return fmt.Sprintf("%s/MultiLanguage/lesson/student/hom_wk_detail.jsp?id=%s&course_id=%s", BaseURL, homeworkId, courseId)
}

func SubmissionURL(courseId, homeworkId string) string {
	return fmt.Sprintf("%s/MultiLanguage/lesson/student/hom_wk_view.jsp?id=%s&course_id=%s", BaseURL, homeworkId, courseId)
}

func HomeworksURL(courseId string) string {
	return fmt.Sprintf("%s/MultiLanguage/lesson/student/hom_wk_brw.jsp?course_id=%s&language=cn", BaseURL, courseId)
}

func (ada *Adapter) HomeworkDetail(courseId, id string, _ map[string]string, homework *resource.Homework) (status int) {
	if homework == nil {
		glog.Errorf("nil received")
		return http.StatusInternalServerError
	}

	url := HomeworkDetailURL(courseId, id)
	doc, err := ada.GetDocument(url)
	if err != nil {
		return http.StatusBadGateway
	}

	// Body.
	bodyTr := doc.Find("table#table_box tr:nth-child(2)")
	if bodyTr.Size() == 0 {
		glog.Errorf("Failed to find the body of the homework.")
		return http.StatusBadGateway
	}
	body := bodyTr.Find("td ~ td").First().Children().Text()

	// Attachment.
	var attach *resource.Attachment
	hrefSelection := bodyTr.Next().Find("td a")
	if fileHref, _ := hrefSelection.Attr("href"); fileHref != "" {
		url := BaseURL + fileHref
		attach = &resource.Attachment{DownloadURL: url}
		if status := ada.FileInfo(url, &attach.Filename, &attach.Size); status != http.StatusOK {
			return status
		}
	}

	// We are safe.
	homework.Body = body
	homework.Attachment = attach
	return http.StatusOK
}

func (ada *Adapter) Submission(courseId, id string, _ map[string]string, submission *resource.Submission) (status int) {
	if submission == nil {
		glog.Errorf("nil received")
		return http.StatusInternalServerError
	}

	url := SubmissionURL(courseId, id)
	doc, err := ada.GetDocument(url)
	if err != nil {
		return http.StatusBadGateway
	}

	tableSelection := doc.Find("#table_box")
	if tableSelection.Length() == 0 {
		glog.Errorf("#table_box not found in submission page")
		return http.StatusBadGateway
	}

	// Body
	infoTr := tableSelection.Find("tr:nth-child(2)")
	submission.Body = infoTr.Find("td.title+td").Children().Text()

	// Attachment
	infoTr = infoTr.Next()
	hrefSelection := infoTr.Find("td a")
	if fileHref, _ := hrefSelection.Attr("href"); fileHref != "" {
		url := BaseURL + fileHref
		attach := &resource.Attachment{DownloadURL: url}
		if status := ada.FileInfo(url, &attach.Filename, &attach.Size); status != http.StatusOK {
			return status
		}
		submission.Attachment = attach
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
		submission.Marked = true
		submission.MarkedAt = markedAt
		if score := strings.TrimSpace(score); strings.Contains(score, "分") {
			if mark, err := strconv.ParseFloat(strings.TrimSuffix(score, "分"), 32); err != nil {
				glog.Errorf("Failed to prase %s into mark: %s", score, err)
				return http.StatusBadGateway
			} else {
				submission.Mark = float32(mark)
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
		url := BaseURL + fileHref
		attach := &resource.Attachment{DownloadURL: url}
		if status := ada.FileInfo(url, &attach.Filename, &attach.Size); status != http.StatusOK {
			return status
		}
		submission.CommentAttachment = attach
	}

	return http.StatusOK
}

func (ada *Adapter) Homeworks(courseId string, _ map[string]string, homeworks *[]*resource.Homework) (status int) {
	if homeworks == nil {
		glog.Errorf("nil received")
		return http.StatusInternalServerError
	}

	doc, err := ada.GetDocument(HomeworksURL(courseId))
	if err != nil {
		return http.StatusBadGateway
	}

	trs := doc.Find("tr.tr1, tr.tr2")
	*homeworks = make([]*resource.Homework, trs.Size())

	statuses := make(chan int, 1)
	count := 0

	trs.Each(func(i int, s *goquery.Selection) {
		infos := s.Find("td").Map(func(i int, tdSelection *goquery.Selection) (info string) {
			info = strings.TrimSpace(tdSelection.Text())
			return info
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
				CourseId: courseId,
				BeginAt:  infos[1],
				DueAt:    infos[2],
				Title:    title,
			}

			// Detile.
			count++
			go func() {
				statuses <- ada.HomeworkDetail(courseId, homeworkId, nil, hw)
			}()
			if infos[3] != "尚未提交" {
				hw.Submissions = append(hw.Submissions, &resource.Submission{HomeworkId: homeworkId})
				count++
				go func() {
					statuses <- ada.Submission(courseId, homeworkId, nil, hw.Submissions[0])
				}()
			}
			(*homeworks)[i] = hw
		}
	})

	status = http.StatusOK
	for i := 0; i < count; i++ {
		if s := <-statuses; s != http.StatusOK {
			status = s
		}
	}

	return status
}
