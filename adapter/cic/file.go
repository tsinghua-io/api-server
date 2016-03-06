package cic

import (
	"encoding/json"
	"fmt"
	"github.com/tsinghua-io/api-server/adapter"
	"github.com/tsinghua-io/api-server/resource"
	"io"
	"strconv"
	"strings"
)

const (
	CoursefilesURL = BaseURL + "/b/myCourse/tree/getCoursewareTreeData/{course_id}/0"
)

type filesParser struct {
	params map[string]string
	data   struct {
		ResultList map[string]struct {
			NodeName     string
			ChildMapData map[string]struct {
				CourseOutlines struct {
					Title string
				}
				CourseCoursewareList []struct {
					ResourcesMappingByFileId struct {
						FileId   string
						RegDate  int64
						FileName string
						FileSize string
						CourseId string
						UserCode string
					}
					RegUser string
					Title   string
					Detail  string
				}
			}
		}
	}
}

func (p *filesParser) Parse(r io.Reader, info interface{}) error {
	files, ok := info.(*[]*resource.File)
	if !ok {
		return fmt.Errorf("The parser and the destination type do not match.")
	}

	dec := json.NewDecoder(r)
	if err := dec.Decode(&p.data); err != nil {
		return err
	}

	for _, node1 := range p.data.ResultList {
		for _, node2 := range node1.ChildMapData {
			category := []string{node1.NodeName, node2.CourseOutlines.Title}

			for _, item := range node2.CourseCoursewareList {
				fileId := item.ResourcesMappingByFileId.FileId
				size, _ := strconv.Atoi(item.ResourcesMappingByFileId.FileSize)

				file := &resource.File{
					Id:       fileId,
					CourseId: item.ResourcesMappingByFileId.CourseId,
					Owner: &resource.User{
						Id:   item.ResourcesMappingByFileId.UserCode,
						Name: item.RegUser,
					},
					CreatedAt:   parseRegDate(item.ResourcesMappingByFileId.RegDate),
					Title:       item.Title,
					Description: item.Detail,
					Category:    category,
					Filename:    item.ResourcesMappingByFileId.FileName,
					Size:        size,
					DownloadUrl: fileID2DownloadUrl(fileId),
				}
				*files = append(*files, file)
			}
		}
	}

	return nil
}

func (ada *CicAdapter) CourseFiles(courseId string, params map[string]string) (files []*resource.File, status int) {
	url := strings.Replace(CoursefilesURL, "{course_id}", courseId, -1)
	parser := &filesParser{params: params}
	files = []*resource.File{}

	status = adapter.FetchInfo(&ada.client, url, "GET", parser, &files)
	return files, status
}
