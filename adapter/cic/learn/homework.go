package learn

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/tsinghua-io/api-server/model"
	"net/http"
	"strconv"
)

func HomeworksURL(courseId string) string {
	return fmt.Sprintf("%s/b/myCourse/homework/list4Student/%s/0", BaseURL, courseId)
}

func (ada *Adapter) Homeworks(courseId string, _ map[string]string, homeworks *[]*model.Homework) (status int) {
	if homeworks == nil {
		glog.Errorf("nil received")
		return http.StatusInternalServerError
	}

	url := HomeworksURL(courseId)
	var v struct {
		ResultList []struct {
			CourseHomeworkRecord struct {
				StudentId                     string
				RegDate                       int64
				HomewkDetail                  string
				ResourcesMappingByHomewkAffix struct {
					FileId   string
					FileName string
					FileSize string
				}
				ReplyDetail                  string
				ResourcesMappingByReplyAffix struct {
					FileId   string
					FileName string
					FileSize string
				}
				Mark      *int
				ReplyDate int64
				Status    string // 0 for 未交, 1 for 未批, 2 for 已阅, 3 for 已批
				IfDelay   string // 1 for late, 2 for 代交
				GradeUser string
			}
			CourseHomeworkInfo struct {
				HomewkId            int
				RegDate             int64
				BeginDate           int64
				EndDate             int64
				Title               string
				Detail              string
				HomewkAffix         string // File ID.
				HomewkAffixFilename string
				CourseId            string
			}
		}
	}

	if err := ada.GetJSON("GET", url, &v); err != nil {
		return http.StatusBadGateway
	}

	for _, result := range v.ResultList {
		// Fetch homework attachment, if exists.
		var attach *model.Attachment

		if fileId := result.CourseHomeworkInfo.HomewkAffix; fileId != "" {
			attach = &model.Attachment{
				Filename:    result.CourseHomeworkInfo.HomewkAffixFilename,
				DownloadURL: DownloadURL(fileId),
			}
			// Get file size.
			if status := ada.FileInfo(attach.DownloadURL, nil, &attach.Size); status != http.StatusOK {
				return status
			}
		}

		// Fetch submission, if exists.
		var submissions []*model.Submission

		if result.CourseHomeworkRecord.Status != "0" {
			// Fetch submission attachment, if exists.
			var attach *model.Attachment
			if affix := result.CourseHomeworkRecord.ResourcesMappingByHomewkAffix; affix.FileId != "" {
				size, _ := strconv.Atoi(affix.FileSize)

				attach = &model.Attachment{
					Filename:    affix.FileName,
					Size:        size,
					DownloadURL: DownloadURL(affix.FileId),
				}
			}

			// Fetch comment attachment, if exists.
			var commentAttach *model.Attachment
			if affix := result.CourseHomeworkRecord.ResourcesMappingByReplyAffix; affix.FileId != "" {
				size, _ := strconv.Atoi(affix.FileSize)

				commentAttach = &model.Attachment{
					Filename:    affix.FileName,
					Size:        size,
					DownloadURL: DownloadURL(affix.FileId),
				}
			}

			// Mark.
			var marked bool
			var mark float32
			if result.CourseHomeworkRecord.Mark != nil {
				marked = true
				mark = float32(*result.CourseHomeworkRecord.Mark)
			}

			submissions = []*model.Submission{
				{
					Owner: &model.User{
						Id: result.CourseHomeworkRecord.StudentId,
					},
					CreatedAt:  parseRegDate(result.CourseHomeworkRecord.RegDate),
					Late:       result.CourseHomeworkRecord.IfDelay == "1",
					Body:       result.CourseHomeworkRecord.HomewkDetail,
					Attachment: attach,
					Marked:     marked,
					MarkedBy: &model.User{
						Name: result.CourseHomeworkRecord.GradeUser,
					},
					MarkedAt:          parseRegDate(result.CourseHomeworkRecord.ReplyDate),
					Mark:              mark,
					Comment:           result.CourseHomeworkRecord.ReplyDetail,
					CommentAttachment: commentAttach,
				},
			}
		}

		homework := &model.Homework{
			Id:          strconv.Itoa(result.CourseHomeworkInfo.HomewkId),
			CourseId:    result.CourseHomeworkInfo.CourseId,
			CreatedAt:   parseRegDate(result.CourseHomeworkInfo.RegDate),
			BeginAt:     parseRegDate(result.CourseHomeworkInfo.BeginDate),
			DueAt:       parseRegDate(result.CourseHomeworkInfo.EndDate),
			Title:       result.CourseHomeworkInfo.Title,
			Body:        result.CourseHomeworkInfo.Detail,
			Attachment:  attach,
			Submissions: submissions,
		}
		*homeworks = append(*homeworks, homework)
	}

	return http.StatusOK
}
