package learn

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/golang/glog"
	"github.com/tsinghua-io/api-server/resource"
	"net/http"
	"net/url"
	"strings"
)

func FilesURL(courseId string) string {
	return fmt.Sprintf("%s/MultiLanguage/lesson/student/download.jsp?course_id=%s", BaseURL, courseId)
}

func (ada *Adapter) Files(courseId string, _ map[string]string, files *[]*resource.File) (status int) {
	if files == nil {
		glog.Errorf("nil received")
		return http.StatusInternalServerError
	}

	doc, err := ada.GetDocument(FilesURL(courseId))
	if err != nil {
		return http.StatusBadGateway
	}

	// Find all categories
	categories := doc.Find("td.textTD").Map(func(i int, s *goquery.Selection) (info string) {
		info = strings.TrimSpace(s.Text())
		return info
	})

	statuses := make(chan int, 1)
	count := 0

	categoryDivs := doc.Find("div.layerbox")
	categoryDivs.Each(func(i int, div *goquery.Selection) {
		category := categories[i]
		trs := div.Find("#table_box tr~tr")
		trs.Each(func(i int, s *goquery.Selection) {
			infos := s.Find("td").Map(func(i int, tdSelection *goquery.Selection) (info string) {
				info = strings.TrimSpace(tdSelection.Text())
				return info
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
				title := strings.TrimSpace(hrefSelection.Text())
				file := &resource.File{
					Id:          fileId,
					CourseId:    courseId,
					CreatedAt:   infos[4],
					Title:       title,
					Description: infos[2],
					Category:    []string{category},
					DownloadURL: BaseURL + href,
				}

				// Get file info.
				count++
				go func() {
					statuses <- ada.FileInfo(file.DownloadURL, &file.Filename, &file.Size)
				}()

				*files = append(*files, file)
			}

		})
	})

	for i := 0; i < count; i++ {
		status = adapter.MergeStatus(status, <-statuses)
	}

	return status
}
