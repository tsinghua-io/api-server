package learn

import (
	"fmt"
	"github.com/tsinghua-io/api-server/model"
	"golang.org/x/text/encoding"
	"net/http"
	"strconv"
)

func AssignmentsURL(courseId string) string {
	return fmt.Sprintf("%s/b/myCourse/homework/list4Student/%s/0", BaseURL, courseId)
}

func (ada *Adapter) Assignments(courseId string) (assignments []*model.Assignment, status int, errMsg error) {
	assignments = make([]*model.Assignment, 0)
	status = http.StatusOK

	url := AssignmentsURL(courseId)
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
				Mark      *float32
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

	if err := ada.GetJSON(url, &v); err != nil {
		return nil, http.StatusBadGateway, err
	}

	for _, result := range v.ResultList {
		id := strconv.Itoa(result.CourseHomeworkInfo.HomewkId)

		// Fetch homework attachment, if exists.
		var attach *model.Attachment

		if fileId := result.CourseHomeworkInfo.HomewkAffix; fileId != "" {
			attach = &model.Attachment{
				Filename:    result.CourseHomeworkInfo.HomewkAffixFilename,
				DownloadURL: DownloadURL(fileId),
			}
			// Get file size.
			if _, attach.Size, status, errMsg = ada.FileInfo(attach.DownloadURL, encoding.Nop); errMsg != nil {
				return
			}
		}

		// Fetch submission, if exists.
		var submission *model.Submission

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

			submission = &model.Submission{
				Owner:             &model.User{Id: result.CourseHomeworkRecord.StudentId},
				AssignmentId:      id,
				CreatedAt:         parseRegDate(result.CourseHomeworkRecord.RegDate),
				Late:              result.CourseHomeworkRecord.IfDelay == "1",
				Body:              result.CourseHomeworkRecord.HomewkDetail,
				Attachment:        attach,
				MarkedBy:          &model.User{Name: result.CourseHomeworkRecord.GradeUser},
				MarkedAt:          parseRegDate(result.CourseHomeworkRecord.ReplyDate),
				Mark:              result.CourseHomeworkRecord.Mark,
				Comment:           result.CourseHomeworkRecord.ReplyDetail,
				CommentAttachment: commentAttach,
			}
		}

		assignment := &model.Assignment{
			Id:         id,
			CourseId:   result.CourseHomeworkInfo.CourseId,
			CreatedAt:  parseRegDate(result.CourseHomeworkInfo.RegDate),
			BeginAt:    parseRegDate(result.CourseHomeworkInfo.BeginDate),
			DueAt:      parseRegDate(result.CourseHomeworkInfo.EndDate),
			Title:      result.CourseHomeworkInfo.Title,
			Body:       result.CourseHomeworkInfo.Detail,
			Attachment: attach,
			Submission: submission,
		}
		assignments = append(assignments, assignment)
	}

	return
}
