package cic

import (
	"encoding/json"
	"fmt"
	"github.com/tsinghua-io/api-server/resource"
	"io"
	"strconv"
	"strings"
)

const (
	FilesURL    = BaseURL + "/b/myCourse/tree/getCoursewareTreeData/{course_id}/0"
	DownloadURL = BaseURL + "/b/resource/downloadFileStream/{file_id}"
)

func fileID2DownloadUrl(fileID string) string {
	return strings.Replace(DownloadURL, "{file_id}", fileID, -1)
}

type filesParser struct {
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

func (p *filesParser) parse(r io.Reader, info interface{}, _ string) error {
	files, ok := info.(*[]*resource.File)
	if !ok {
		return fmt.Errorf("The parser and the destination type do not match.")
	}

	dec := json.NewDecoder(r)
	if err := dec.Decode(p); err != nil {
		return err
	}

	for _, node1 := range p.ResultList {
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

func (adapter *CicAdapter) Files(course_id string) (files []*resource.File, status int) {
	files = []*resource.File{}
	url := strings.Replace(FilesURL, "{course_id}", course_id, -1)
	status = adapter.FetchInfo(url, "GET", &filesParser{}, &files)
	return files, status
}
