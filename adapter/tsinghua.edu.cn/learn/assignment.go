package learn

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/tsinghua-io/api-server/model"
	"github.com/tsinghua-io/api-server/util"
	"golang.org/x/text/encoding/simplifiedchinese"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func AssignmentDetailURL(courseId, id string) string {
	return fmt.Sprintf("%s/MultiLanguage/lesson/student/hom_wk_detail.jsp?course_id=%s&id=%s", BaseURL, courseId, id)
}

func SubmissionURL(courseId, id string) string {
	return fmt.Sprintf("%s/MultiLanguage/lesson/student/hom_wk_view.jsp?course_id=%s&id=%s&language=cn", BaseURL, courseId, id)
}

func AssignmentsURL(courseId string) string {
	return fmt.Sprintf("%s/MultiLanguage/lesson/student/hom_wk_brw.jsp?course_id=%s&language=cn", BaseURL, courseId)
}

func ParseAssignmentDetailURL(rawurl string) (courseId, id string) {
	if parsed, err := url.Parse(rawurl); err == nil {
		values := parsed.Query()
		courseId = values.Get("course_id")
		id = values.Get("id")
	}
	return
}

func ParseMark(markStr string) (markPtr *float32, err error) {
	if markStr != "" {
		var mark float64
		if mark, err = strconv.ParseFloat(strings.TrimSuffix(markStr, "分"), 32); err == nil {
			markPtr = util.NewFloat32(float32(mark))
		}
	}
	return
}

func (ada *Adapter) Attachment(s *goquery.Selection) (attach *model.Attachment, status int, errMsg error) {
	status = http.StatusOK

	if href, ok := s.Find("a").Attr("href"); ok {
		attach = &model.Attachment{DownloadURL: BaseURL + href}
		if attach.Filename, attach.Size, status, errMsg = ada.FileInfo(attach.DownloadURL, simplifiedchinese.GBK); errMsg != nil {
			errMsg = fmt.Errorf("Failed to get attachment info from %s: %s", href, errMsg)
		}
	}

	return
}

func (ada *Adapter) AssignmentDetail(courseId, id string) (title, body string, attach *model.Attachment, status int, errMsg error) {
	status = http.StatusOK

	url := AssignmentDetailURL(courseId, id)
	doc, err := ada.GetDocument(url)
	if err != nil {
		status = http.StatusBadGateway
		errMsg = err
		return
	}

	var attachSel *goquery.Selection

	// Careful, I remembered assignments could contain HTML before.
	if content := doc.Find("#table_box>tbody>tr>td.title+td"); content.Size() != 5 {
		errMsg = fmt.Errorf("Expect 5 content blocks, got %d.", content.Size())
	} else {
		attachSel = content.Eq(2)
		texts := util.TrimmedTexts(content)
		title = texts[0]
		body = texts[1]
	}

	if errMsg != nil {
		status = http.StatusInternalServerError
		errMsg = fmt.Errorf("Failed to parse %s: %s", url, errMsg)
		return
	}

	attach, status, errMsg = ada.Attachment(attachSel)
	return
}

func (ada *Adapter) Submission(courseId, id string) (submission *model.Submission, status int, errMsg error) {
	status = http.StatusOK

	url := SubmissionURL(courseId, id)
	doc, err := ada.GetDocument(url)
	if err != nil {
		status = http.StatusBadGateway
		errMsg = err
		return
	}

	var attachSel, commentAttachSel *goquery.Selection

	// Careful, I remembered assignments could contain HTML before.
	if content := doc.Find("#table_box>tbody>tr>td.title+td"); content.Size() != 8 {
		errMsg = fmt.Errorf("Expect 8 content blocks, got %d.", content.Size())
	} else {
		texts := util.TrimmedTexts(content)
		markStr := texts[5]
		if mark, err := ParseMark(markStr); err != nil {
			errMsg = fmt.Errorf("Failed to parse mark from %s: %s", markStr, err)
		} else {
			attachSel = content.Eq(2)
			commentAttachSel = content.Eq(7)

			submission = &model.Submission{
				AssignmentId: id,
				Late:         false, // Or we cannot view it now.
				Body:         texts[1],
			}
			if markedAt := texts[4]; markedAt != "null" {
				if name := texts[3]; name != "" {
					submission.MarkedBy = &model.User{Name: name}
				}
				submission.MarkedAt = markedAt
				submission.Mark = mark
				submission.Comment = texts[6]
			}
		}
	}

	if errMsg != nil {
		status = http.StatusInternalServerError
		errMsg = fmt.Errorf("Failed to parse %s: %s", url, errMsg)
		return
	}

	sg := util.NewStatusGroup()
	sg.Go(func(status *int, err *error) {
		submission.Attachment, *status, *err = ada.Attachment(attachSel)
	})
	sg.Go(func(status *int, err *error) {
		submission.CommentAttachment, *status, *err = ada.Attachment(commentAttachSel)
	})

	status, errMsg = sg.Wait()
	return
}

func (ada *Adapter) AssignmentList(courseId string) (assignments []*model.Assignment, status int, errMsg error) {
	status = http.StatusOK

	url := AssignmentsURL(courseId)
	doc, err := ada.GetDocument(url)
	if err != nil {
		status = http.StatusBadGateway
		errMsg = err
		return
	}

	rows := doc.Find("tr.tr1, tr.tr2")
	assignments = make([]*model.Assignment, rows.Size())

	rows.EachWithBreak(func(i int, s *goquery.Selection) bool {
		if cols := s.Children(); cols.Size() != 6 {
			errMsg = fmt.Errorf("Expect 6 columns, got %d.", cols.Size())
		} else if link := cols.Eq(0).Find("a"); link.Size() == 0 {
			errMsg = fmt.Errorf("Failed to find link in column 0.")
		} else if _, id := ParseAssignmentDetailURL(link.AttrOr("href", "")); id == "" {
			errMsg = fmt.Errorf("Failed to find assignment id from href (%s).", link.AttrOr("href", ""))
		} else {
			texts := util.TrimmedTexts(cols)
			assign := &model.Assignment{
				Id:       id,
				CourseId: courseId,
				BeginAt:  texts[1],
				DueAt:    texts[2] + "T23:59:59+0800",
				Title:    texts[0],
			}
			switch subStatus := texts[3]; subStatus {
			case "尚未提交":
			case "已经提交":
				assign.Submission = new(model.Submission)
			default:
				errMsg = fmt.Errorf("Unknown submission status: %s", subStatus)
			}

			assignments[i] = assign
		}

		return errMsg == nil
	})

	if errMsg != nil {
		status = http.StatusInternalServerError
		errMsg = fmt.Errorf("Failed to parse %s: %s", url, errMsg)
		return
	}

	return
}

func (ada *Adapter) Assignments(courseId string) (assignments []*model.Assignment, status int, errMsg error) {
	if assignments, status, errMsg = ada.AssignmentList(courseId); errMsg != nil {
		return
	}

	sg := util.NewStatusGroup()

	for _, assign := range assignments {
		assign := assign

		sg.Go(func(status *int, err *error) {
			_, assign.Body, assign.Attachment, *status, *err = ada.AssignmentDetail(courseId, assign.Id)
		})

		if assign.Submission != nil {
			sg.Go(func(status *int, err *error) {
				assign.Submission, *status, *err = ada.Submission(courseId, assign.Id)
			})
		}
	}

	status, errMsg = sg.Wait()
	return
}
