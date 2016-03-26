package learn

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/tsinghua-io/api-server/resource"
	"net/http"
	"strconv"
	"strings"
)

func FilesURL(courseId string) string {
	return fmt.Sprintf("%s/b/myCourse/tree/getCoursewareTreeData/%s/0", BaseURL, courseId)
}

func DownloadURL(fileId string) string {
	return fmt.Sprintf("%s/b/resource/downloadFileStream/%s", BaseURL, fileId)
}

func (ada *Adapter) Files(courseId string, _ map[string]string, files *[]*resource.File) (status int) {
	if files == nil {
		glog.Errorf("nil received")
		return http.StatusInternalServerError
	}

	url := FilesURL(courseId)
	var v struct {
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

	if err := ada.GetJSON("POST", url, &v); err != nil {
		return http.StatusBadGateway
	}

	for _, node1 := range v.ResultList {
		for _, node2 := range node1.ChildMapData {
			category := []string{
				strings.TrimSpace(node1.NodeName),
				strings.TrimSpace(node2.CourseOutlines.Title),
			}

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
					DownloadURL: DownloadURL(fileId),
				}
				*files = append(*files, file)
			}
		}
	}

	return http.StatusOK
}
