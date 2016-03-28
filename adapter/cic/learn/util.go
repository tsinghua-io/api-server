package learn

import (
	"fmt"
	"net/http"
	"time"
)

const (
	TeachingWeekURL = BaseURL + "/b/myCourse/courseList/getCurrentTeachingWeek"
)

func DownloadURL(fileId string) string {
	return fmt.Sprintf("%s/b/resource/downloadFileStream/%s", BaseURL, fileId)
}

func (ada *Adapter) Semesters() (thisSem, nextSem string, status int) {
	var v struct {
		CurrentSemester struct {
			Id string
		}
		NextSemester struct {
			Id string
		}
	}

	if err := ada.GetJSON("GET", TeachingWeekURL, &v); err != nil {
		status = http.StatusBadGateway
		return
	}

	thisSem = v.CurrentSemester.Id
	nextSem = v.NextSemester.Id
	status = http.StatusOK
	return
}

func parseRegDate(regDate int64) string {
	// Return empty string for 0.
	// Will not work for 1970-01-01T08:00:00+0800.
	if regDate == 0 {
		return ""
	}
	return time.Unix(regDate/1000, 0).Format("2006-01-02T15:04:05+0800")
}
