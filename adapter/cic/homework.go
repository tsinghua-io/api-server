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
	HomeworksURL = BaseURL + "/b/myCourse/homework/list4Student/{course_id}/0"
)

func Int2Pfloat32(i int) *float32 {
	f := float32(i)
	return &f
}

type homeworksParser struct {
	ResultList []struct {
		CourseHomeworkRecord struct {
			StudentId                     string
			RegDate                       int64
			HomewkDetail                  string
			ResourcesMappingByHomewkAffix struct {
				FileId   string
				FileName string
				FileSize string
				UserCode string
			}
			ReplyDetail string
			// TODO: Add this:
			// ResourcesMappingByReplyAffix struct {
			// }
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
			// AnswerDetail
			// AnswerLink
			// AnswerLinkFilename
			// AnswerDate
			CourseId string
			WeiJiao  int
			// YiJiao
			YiYue  int
			YiPi   int
			Jiaoed int
		}
	}
}

func (p *homeworksParser) parse(r io.Reader, info interface{}, _ string) error {
	homeworks, ok := info.(*[]*resource.Homework)
	if !ok {
		return fmt.Errorf("The parser and the destination type do not match.")
	}

	dec := json.NewDecoder(r)
	if err := dec.Decode(p); err != nil {
		return err
	}

	for _, result := range p.ResultList {
		// Fetch homework attachment, if exists.
		var attach *resource.Attachment

		if fileID := result.CourseHomeworkInfo.HomewkAffix; fileID != "" {
			attach = &resource.Attachment{
				Filename: result.CourseHomeworkInfo.HomewkAffixFilename,
				// TODO: Get file size?
				DownloadUrl: fileID2DownloadUrl(fileID),
			}
		}

		// Fetch submission, if exists.
		var submissions []*resource.Submission

		if result.CourseHomeworkRecord.Status != "0" {
			// Fetch submission attachment, if exists.
			var attach *resource.Attachment

			if affix := result.CourseHomeworkRecord.ResourcesMappingByHomewkAffix; affix.FileId != "" {
				size, _ := strconv.Atoi(affix.FileSize)

				attach = &resource.Attachment{
					Filename:    affix.FileName,
					Size:        size,
					DownloadUrl: fileID2DownloadUrl(affix.FileId),
				}
			}

			// Fetch mark, if exists.
			var mark *float32
			if intMark := result.CourseHomeworkRecord.Mark; intMark != nil {
				mark = Int2Pfloat32(*intMark)
			}

			submissions = []*resource.Submission{
				{
					Owner: &resource.User{
						Id: result.CourseHomeworkRecord.StudentId,
					},
					CreatedAt:  parseRegDate(result.CourseHomeworkRecord.RegDate),
					Late:       result.CourseHomeworkRecord.IfDelay == "1",
					Body:       result.CourseHomeworkRecord.HomewkDetail,
					Attachment: attach,
					Mark:       mark,
					MarkedBy: &resource.User{
						Name: result.CourseHomeworkRecord.GradeUser,
					},
					MarkedAt: parseRegDate(result.CourseHomeworkRecord.ReplyDate),
					Comment:  result.CourseHomeworkRecord.ReplyDetail,
					// TODO: Add this.
					// CommentAttachment: resource.Attachment{
					// }
				},
			}
		}

		homework := &resource.Homework{
			Id:                strconv.Itoa(result.CourseHomeworkInfo.HomewkId),
			CourseId:          result.CourseHomeworkInfo.CourseId,
			CreatedAt:         parseRegDate(result.CourseHomeworkInfo.RegDate),
			BeginAt:           parseRegDate(result.CourseHomeworkInfo.BeginDate),
			DueAt:             parseRegDate(result.CourseHomeworkInfo.EndDate),
			SubmittedCount:    result.CourseHomeworkInfo.Jiaoed,
			NotSubmittedCount: result.CourseHomeworkInfo.WeiJiao,
			SeenCount:         result.CourseHomeworkInfo.YiYue,
			MarkedCount:       result.CourseHomeworkInfo.YiPi,
			Title:             result.CourseHomeworkInfo.Title,
			Body:              result.CourseHomeworkInfo.Detail,
			Attachment:        attach,
			Submissions:       submissions,
		}
		*homeworks = append(*homeworks, homework)
	}

	return nil
}

func (adapter *CicAdapter) Homeworks(course_id string) (homeworks []*resource.Homework, status int) {
	homeworks = []*resource.Homework{}
	url := strings.Replace(HomeworksURL, "{course_id}", course_id, -1)
	status = adapter.FetchInfo(url, "GET", &homeworksParser{}, &homeworks)
	return homeworks, status
}
