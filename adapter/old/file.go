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
	"strings"
)

const (
	courseFileURL = BaseURL + "/MultiLanguage/lesson/student/download.jsp?course_id={course_id}"
)

type courseFilesParser struct {
	params   map[string]string
	courseId string
}

func (p *courseFilesParser) Parse(r io.Reader, info interface{}) error {
	files, ok := info.(*[]*resource.File)
	if !ok {
		return fmt.Errorf("The parser and the destination type do not match.")
	}

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return err
	}

	// Find all categories
	categories := doc.Find("td.textTD").Map(func(i int, s *goquery.Selection) (info string) {
		info = strings.TrimSpace(s.Text())
		return
	})

	categoryDivs := doc.Find("div.layerbox")
	categoryDivs.Each(func(i int, div *goquery.Selection) {
		category := categories[i]
		trs := div.Find("#table_box tr~tr")
		trs.Each(func(i int, s *goquery.Selection) {
			infos := s.Find("td").Map(func(i int, tdSelection *goquery.Selection) (info string) {
				info = strings.TrimSpace(tdSelection.Text())
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
				title := strings.TrimSpace(hrefSelection.Text())
				file := &resource.File{
					Id:          fileId,
					CourseId:    p.courseId,
					CreatedAt:   infos[4],
					Title:       title,
					Description: infos[2],
					Category:    []string{category},
					DownloadUrl: BaseURL + href,
				}

				*files = append(*files, file)
			}
		})
	})

	return nil
}

func (ada *OldAdapter) CourseFiles(courseId string, params map[string]string) (files []*resource.File, status int) {
	URL := strings.Replace(courseFileURL, "{course_id}", courseId, -1)
	parser := &courseFilesParser{params: params, courseId: courseId}

	if status = adapter.FetchInfo(&ada.client, URL, "GET", parser, &files); status != http.StatusOK {
		return nil, status
	}
	count := len(files)
	statuses := make(chan int, count)

	for _, file := range files {
		file := file

		go func() {
			filename, size, err := parseFileInfo(ada.client, file.DownloadUrl)
			if err != nil {
				glog.Errorf("Failed to parse file info of %s: %s", file.DownloadUrl, err)
				statuses <- http.StatusBadGateway
			}
			file.Filename = filename
			file.Size = size
			statuses <- http.StatusOK
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
