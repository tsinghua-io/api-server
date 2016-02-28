package cic

import (
	"encoding/json"
	"fmt"
	"github.com/tsinghua-io/api-server/resource"
	"io"
	"strconv"
)

type homeworksParser struct {
	ResultList []struct {
		courseHomeworkRecord struct {
			studentId                     string
			teacherId                     string
			regDate                       int64
			homewkDetail                  string
			resourcesMappingByHomewkAffix struct {
				fileId   string
				regDate  string
				fileName string
				fileSize string
				userCode string
			}
			replyDetail string
			// TODO: Add this:
			// resourcesMappingByReplyAffix struct {
			// }
			mark      int
			replyDate int64
			status    string // 0 for 未交, 1 for 未批, 2 for 已阅, 3 for 已批
			ifDelay   string // 1 for late, 2 for 代交
			gradeUser string
		}
		courseHomeworkInfo struct {
			homewkId            int
			regDate             int64
			beginDate           int64
			endDate             int64
			title               string
			detail              string
			homewkAffix         string // File ID.
			homewkAffixFilename string
			// answerDetail
			// answerLink
			// answerLinkFilename
			// answerDate
			courseId string
			weiJiao  int
			// yiJiao
			yiYue  int
			yiPi   int
			jiaoed int
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
		// Fetch attachments, if existed.
		var attach, submissionAttach *resource.Attachment

		if fileID := result.courseHomeworkInfo.homewkAffix; fileID != "" {
			attach = &resource.Attachment{
				Filename: result.courseHomeworkInfo.homewkAffixFilename,
				// TODO: Get file size?
				DownloadUrl: fileID2DownloadUrl(fileID),
			}
		}

		if status := result.courseHomeworkRecord.status; status == "0" {
			// Not submitted.

		}
		if fileID := result.courseHomeworkRecord.resourcesMappingByHomewkAffix.fileId; fileID != "" {
			size, _ := strconv.Atoi(result.courseHomeworkRecord.resourcesMappingByHomewkAffix.fileSize)

			submissionAttach = &resource.Attachment{
				Filename:    result.courseHomeworkRecord.resourcesMappingByHomewkAffix.fileName,
				Size:        size,
				DownloadUrl: fileID2DownloadUrl(fileID),
			}
		}

		mark := float32(result.courseHomeworkRecord.mark)

		homework := &resource.Homework{
			Id:                strconv.Itoa(result.courseHomeworkInfo.homewkId),
			CourseId:          result.courseHomeworkInfo.courseId,
			CreatedAt:         parseRegDate(result.courseHomeworkInfo.regDate),
			BeginAt:           parseRegDate(result.courseHomeworkInfo.beginDate),
			DueAt:             parseRegDate(result.courseHomeworkInfo.endDate),
			SubmittedCount:    result.courseHomeworkInfo.jiaoed,
			NotSubmittedCount: result.courseHomeworkInfo.weiJiao,
			SeenCount:         result.courseHomeworkInfo.yiYue,
			MarkedCount:       result.courseHomeworkInfo.yiPi,
			Title:             result.courseHomeworkInfo.title,
			Body:              result.courseHomeworkInfo.detail,
			Attachment:        attach,
			Submissions: []*resource.Submission{
				&resource.Submission{
					Owner: &resource.User{
						Id: result.courseHomeworkRecord.studentId,
					},
					CreatedAt:  parseRegDate(result.courseHomeworkRecord.regDate),
					Late:       result.courseHomeworkRecord.ifDelay == "1",
					Body:       result.courseHomeworkRecord.homewkDetail,
					Attachment: submissionAttach,
					Mark:       &mark,
					MarkedBy: &resource.User{
						Id:   result.courseHomeworkRecord.teacherId,
						Name: result.courseHomeworkRecord.gradeUser,
					},
					MarkedAt: parseRegDate(result.courseHomeworkRecord.replyDate),
					Comment:  result.courseHomeworkRecord.replyDetail,
					// TODO: Add this.
					// CommentAttachment: resource.Attachment{
					// }
				},
			},
		}
		*homeworks = append(*homeworks, homework)
	}

	return nil
}

func (adapter *CicAdapter) Homeworks(course_id string) (courses []*resource.Homework, status int) {
	return
}
