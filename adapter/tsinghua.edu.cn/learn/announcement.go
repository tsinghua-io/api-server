package learn

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/tsinghua-io/api-server/model"
	"github.com/tsinghua-io/api-server/util"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func AnnouncementURL(courseId, id string) string {
	return fmt.Sprintf("%s/MultiLanguage/public/bbs/note_reply.jsp?bbs_type=课程公告&course_id=%s&id=%s", BaseURL, courseId, id)
}

func ParseAnnouncementURL(rawurl string) (courseId, id string) {
	if parsed, err := url.Parse(rawurl); err == nil {
		values := parsed.Query()
		courseId = values.Get("course_id")
		id = values.Get("id")
	}
	return
}

func AnnouncementListURL(courseId string) string {
	return fmt.Sprintf("%s/MultiLanguage/public/bbs/getnoteid_student.jsp?course_id=%s", BaseURL, courseId)
}

func (ada *Adapter) Announcement(courseId, id string) (title, body, email string, status int, errMsg error) {
	status = http.StatusOK

	url := AnnouncementURL(courseId, id)
	doc, err := ada.GetDocument(url)
	if err != nil {
		status = http.StatusBadGateway
		errMsg = err
		return
	}

	// Be careful here, as the body can contain anything.
	if content := doc.Find("#table_box>tbody>tr>td:last-child"); content.Size() != 2 {
		errMsg = fmt.Errorf("Expect 2 content blocks, got %d.", content.Size())
	} else if bodyHTML, err := content.Eq(1).Html(); err != nil {
		errMsg = fmt.Errorf("Failed to rebuild body HTML: %s", err)
	} else {
		title = strings.TrimSpace(content.Eq(0).Text())
		body = strings.TrimSpace(bodyHTML)
		if sendEmail := doc.Find("#table_box>tbody>tr:last-child>td>input[name=sendmail]"); sendEmail.Size() != 0 {
			// Try our best to parse it, but don't panic if we cannot.
			regex, _ := regexp.Compile("/MultiLanguage/public/mail/student/sendmail.jsp\\?usersToSend=(.*?)&")
			if subs := regex.FindStringSubmatch(sendEmail.AttrOr("onclick", "")); len(subs) > 1 {
				email = subs[1]
			}
		}
	}

	if errMsg != nil {
		status = http.StatusInternalServerError
		errMsg = fmt.Errorf("Failed to parse %s: %s", url, errMsg)
		return
	}

	return
}

func (ada *Adapter) AnnouncementList(courseId string) (announcements []*model.Announcement, status int, errMsg error) {
	status = http.StatusOK

	url := AnnouncementListURL(courseId)
	doc, err := ada.GetDocument(url)
	if err != nil {
		status = http.StatusBadGateway
		errMsg = err
		return
	}

	rows := doc.Find("#table_box tr~tr")
	announcements = make([]*model.Announcement, rows.Size())

	rows.EachWithBreak(func(i int, s *goquery.Selection) bool {
		if cols := s.Children(); cols.Size() != 5 {
			errMsg = fmt.Errorf("Expect 5 columns, got %d.", cols.Size())
		} else if link := cols.Eq(1).Find("a"); link.Size() == 0 {
			errMsg = fmt.Errorf("Failed to find link in column 1.")
		} else if _, id := ParseAnnouncementURL(link.AttrOr("href", "")); id == "" {
			errMsg = fmt.Errorf("Failed to find announcement id from href (%s).", link.AttrOr("href", ""))
		} else {
			texts := util.TrimmedTexts(cols)
			annc := &model.Announcement{
				Id:        id,
				CourseId:  courseId,
				Owner:     &model.User{Name: strings.TrimSuffix(texts[2], "老师")},
				CreatedAt: texts[3],
				Title:     texts[1],
			}
			if red := link.Find("font[color=red]"); red.Size() > 0 {
				annc.Priority = 1
			}

			announcements[i] = annc
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

func (ada *Adapter) Announcements(courseId string) (announcements []*model.Announcement, status int, errMsg error) {
	if announcements, status, errMsg = ada.AnnouncementList(courseId); errMsg != nil {
		return
	}

	sg := util.NewStatusGroup()
	sg.Add(len(announcements))

	for _, annc := range announcements {
		annc := annc
		go func() {
			var status int
			var err error
			defer sg.Done(status, err)
			_, annc.Body, annc.Owner.Email, status, err = ada.Announcement(courseId, annc.Id)
		}()
	}

	status, errMsg = sg.Wait()
	return
}
