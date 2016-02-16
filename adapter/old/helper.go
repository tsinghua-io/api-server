package old

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"github.com/golang/glog"
	"github.com/tsinghua-io/api-server/resource"
	"golang.org/x/net/html"
	"mime"
	"net/url"
	"strconv"
	"strings"
)

func (adapter *OldAdapter) courseIds(typepage int) (courseIdList []string, err error) {
	path := "/MultiLanguage/lesson/student/MyCourse.jsp?typepage=" + strconv.Itoa(typepage)
	doc, err := adapter.getOldResponse(path, make(map[string]string))

	if err != nil {
		err = fmt.Errorf("Failed to get response from learning web: %s", err)
	} else {
		// parsing the response body
		courseLinkList := doc.Find("#info_1 tr a")
		courseLinkList.Each(func(i int, s *goquery.Selection) {
			var href string
			var hrefUrl *url.URL
			var courseId string
			href, _ = s.Attr("href")
			if hrefUrl, err = url.Parse(href); err != nil {
				return
			}
			if courseId = hrefUrl.Query().Get("course_id"); courseId != "" {
				courseIdList = append(courseIdList, courseId)
			}
		})
	}
	return
}

func (adapter *OldAdapter) courseInfo(courseId string) (course *resource.Course, err error) {
	path := "/MultiLanguage/lesson/student/course_info.jsp?course_id=" + courseId
	doc, err := adapter.getOldResponse(path, make(map[string]string))

	if err != nil {
		err = fmt.Errorf("Failed to get response from learning web: %s", err)
	} else {
		tds := doc.Find("table#table_box td")

		infos := tds.Map(func(i int, s *goquery.Selection) (info string) {
			firstChild := s.Nodes[0].FirstChild
			if firstChild != nil {
				switch s.Nodes[0].FirstChild.Type {
				case html.TextNode:
					info, _ = s.Html()
				case html.ElementNode:
					info, _ = s.Children().Html()
				default:
					info, _ = s.Html()
				}
			}
			return
		})
		if len(infos) < 23 {
			err = fmt.Errorf("Course information parsing error: cannot parse all the informations from %s", infos)
			return
		}
		course = &resource.Course{
			Id:   courseId,
			Name: infos[5],
			Teacher: resource.User{
				Name:  infos[16],
				Email: infos[18],
				Phone: infos[20],
			},
			CourseNumber:   infos[1],
			CourseSequence: infos[3],
			Description:    infos[22],
		}

		if credit, err := strconv.Atoi(strings.TrimSpace(infos[7])); err == nil {
			course.Credit = credit
		}
		if hour, err := strconv.Atoi(strings.TrimSpace(infos[9])); err == nil {
			course.Hour = hour
		}
	}
	return
}

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

func (adapter *OldAdapter) parseFileInfo(path string) (filename string, size int) {
	resp, err := adapter.client.Head(BaseURL + path)
	if err != nil {
		glog.Errorf("Failed to get header information of file %s: %s", path, err)
		return
	}

	// file size
	size, _ = strconv.Atoi(resp.Header.Get("Content-Length"))

	// file name
	disposition := resp.Header.Get("Content-Disposition")
	// decode from gbk
	disposition = mahonia.NewDecoder("GBK").ConvertString(disposition)

	// parse disposition header
	disposition, params, err := mime.ParseMediaType(disposition)
	if err != nil {
		glog.Errorf("Failed to parse header Content-Disposition of file %s", path)
		return
	}
	filename = params["filename"]
	return
}

func (adapter *OldAdapter) parseHomeworkInfo(href string) (body string, attachment resource.Attachment) {
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
			filename, size := adapter.parseFileInfo(fileHref)
			attachment = resource.Attachment{
				Filename:    filename,
				Size:        size,
				DownloadUrl: href,
			}
		}

	}
	return
}
