package learn

import (
	"fmt"
	"github.com/tsinghua-io/api-server/model"
	"net/http"
	"time"
)

const (
	TeachingWeekURL = BaseURL + "/b/myCourse/courseList/getCurrentTeachingWeek"
)

func DownloadURL(fileId string) string {
	return fmt.Sprintf("%s/b/resource/downloadFileStream/%s", BaseURL, fileId)
}

func parseRegDate(regDate int64) string {
	// Return empty string for 0.
	// Will not work for 1970-01-01T08:00:00+0800.
	if regDate == 0 {
		return ""
	}
	return time.Unix(regDate/1000, 0).Format("2006-01-02T15:04:05+0800")
}

func (ada *Adapter) Semesters() (thisSem, nextSem *model.Semester, status int, errMsg error) {
	status = http.StatusOK

	var v struct {
		CurrentSemester struct {
			Id        string
			StartDate string
		}
		NextSemester struct {
			Id        string
			StartDate string
		}
	}

	if err := ada.GetJSON(TeachingWeekURL, &v); err != nil {
		status = http.StatusBadGateway
		errMsg = err
		return
	}

	thisSem = &model.Semester{Id: v.CurrentSemester.Id, BeginAt: v.CurrentSemester.StartDate}
	nextSem = &model.Semester{Id: v.NextSemester.Id, BeginAt: v.NextSemester.StartDate}
	return
}
