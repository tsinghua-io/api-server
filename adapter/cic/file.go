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
	downloadURL = "http://learn.cic.tsinghua.edu.cn/b/resource/downloadFileStream/{file_id}"
)

func fileID2DownloadUrl(fileID string) string {
	return strings.Replace(downloadURL, "{file_id}", fileID, -1)
}

type filesParser struct {
	ResultList map[string]struct {
		TeacherInfoView struct {
			NodeName     string
			ChildMapData map[string]struct {
				Title                string
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
		for _, node2 := range node1.TeacherInfoView.ChildMapData {
			category := []string{node1.TeacherInfoView.NodeName, node2.Title}

			for _, item := range node2.CourseCoursewareList {
				fileId := item.ResourcesMappingByFileId.FileId
				size, _ := strconv.Atoi(item.ResourcesMappingByFileId.FileSize)

				file := &resource.File{
					Id:       fileId,
					CourseId: item.ResourcesMappingByFileId.CourseId,
					Owner: resource.User{
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

func (adapter *CicAdapter) Files(course_id string) (courses []*resource.File, status int) {
	return
}
