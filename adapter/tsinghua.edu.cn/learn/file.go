package learn

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/tsinghua-io/api-server/model"
	"github.com/tsinghua-io/api-server/util"
	"golang.org/x/text/encoding/simplifiedchinese"
	"net/http"
	"net/url"
)

func FilesURL(courseId string) string {
	return fmt.Sprintf("%s/MultiLanguage/lesson/student/download.jsp?course_id=%s", BaseURL, courseId)
}

func ParseDownloadURL(downloadURL string) (courseId, id string) {
	if parsed, err := url.Parse(downloadURL); err == nil {
		courseId = parsed.Query().Get("course_id")
		id = parsed.Query().Get("file_id")
	}
	return
}

func (ada *Adapter) Files(courseId string) (files []*model.File, status int, errMsg error) {
	// TODO: Clean this function up.
	files = make([]*model.File, 0)

	url := FilesURL(courseId)
	doc, err := ada.GetDocument(url)
	if err != nil {
		status = http.StatusBadGateway
		errMsg = err
		return
	}

	// Find all categories.
	categories := util.TrimmedTexts(doc.Find("#table_box td.textTD"))

	doc.Find("div.layerbox").EachWithBreak(func(i int, div *goquery.Selection) bool {
		if i >= len(categories) {
			return false
		}
		category := []string{categories[i]}

		div.Find("#table_box tr~tr").EachWithBreak(func(_ int, s *goquery.Selection) bool {
			if cols := s.Children(); cols.Size() != 6 {
				errMsg = fmt.Errorf("Expect 6 columns, got %d.", cols.Size())
			} else if link := cols.Eq(1).Find("a"); link.Size() == 0 {
				errMsg = fmt.Errorf("Failed to find link in column 1.")
			} else if _, id := ParseDownloadURL(link.AttrOr("href", "")); id == "" {
				errMsg = fmt.Errorf("Failed to find file id from href (%s).", link.AttrOr("href", ""))
			} else {
				texts := util.TrimmedTexts(cols)
				file := &model.File{
					Id:          id,
					CourseId:    courseId,
					CreatedAt:   texts[4],
					Title:       texts[1],
					Description: texts[2],
					Category:    category,
					DownloadURL: BaseURL + link.AttrOr("href", ""),
				}

				files = append(files, file)
			}

			return errMsg == nil
		})

		return errMsg == nil
	})

	if errMsg != nil {
		status = http.StatusInternalServerError
		errMsg = fmt.Errorf("Failed to parse %s: %s", url, errMsg)
		return
	}

	// Fill file infos.
	sg := util.NewStatusGroup()

	for _, file := range files {
		file := file
		sg.Go(func(status *int, err *error) {
			file.Filename, file.Size, *status, *err = ada.FileInfo(file.DownloadURL, simplifiedchinese.GBK)
		})
	}

	status, errMsg = sg.Wait()
	return
}
