package old

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/golang/glog"
	"github.com/tsinghua-io/api-server/resource"
	"net/http"
	"net/url"
	"strings"
)

func (adapter *OldAdapter) Files(courseId string) (files []*resource.File, status int) {
	path := "/MultiLanguage/lesson/student/download.jsp?course_id=" + courseId
	doc, err := adapter.getOldResponse(path, make(map[string]string))

	status = http.StatusBadGateway

	if err != nil {
		glog.Errorf("Failed to get response from learning web: %s", err)
	} else {
		// Find all categories
		categories := doc.Find("td.textTD").Map(func(i int, s *goquery.Selection) (info string) {
			info, _ = s.Html()
			return
		})

		categoryDivs := doc.Find("div.layerbox")
		categoryDivs.Each(func(i int, div *goquery.Selection) {
			category := categories[i]
			trs := div.Find("#table_box tr~tr")
			trs.Each(func(i int, s *goquery.Selection) {
				infos := s.Find("td").Map(func(i int, tdSelection *goquery.Selection) (info string) {
					info, _ = tdSelection.Html()
					return
				})

				hrefSelection := s.Find("td a")

				var href string
				var hrefUrl *url.URL
				var err error
				if href, _ = hrefSelection.Attr("href"); href == "" {
					return
				}
				if hrefUrl, err = url.Parse(href); err != nil {
					return
				}

				if fileId := hrefUrl.Query().Get("file_id"); fileId != "" {
					title, _ := hrefSelection.Html()
					title = strings.TrimSpace(title)
					file := &resource.File{
						Id:          fileId,
						CourseId:    courseId,
						Category:    []string{category},
						Title:       title,
						Description: infos[2],
						DownloadUrl: href,
						CreatedAt:   infos[4],
					}

					file.Filename, file.Size = adapter.parseFileInfo(href)
					files = append(files, file)
				}
			})
		})
		status = http.StatusOK
	}
	return
}
