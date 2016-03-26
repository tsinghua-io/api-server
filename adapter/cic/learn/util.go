package learn

import (
	"net/http"
	"time"
)

const (
	TeachingWeekURL = BaseURL + "/b/myCourse/courseList/getCurrentTeachingWeek"
)

func (ada *Adapter) Semesters(currentSemester, nextSemester *string) (status int) {
	var v struct {
		CurrentSemester struct {
			Id string
		}
		NextSemester struct {
			Id string
		}
	}

	if err := ada.GetJSON("GET", TeachingWeekURL, &v); err != nil {
		return http.StatusBadGateway
	}

	if currentSemester != nil {
		*currentSemester = v.CurrentSemester.Id
	}
	if nextSemester != nil {
		*nextSemester = v.NextSemester.Id
	}

	return http.StatusOK
}

func parseRegDate(regDate int64) string {
	// Return empty string for 0.
	// Will not work for 1970-01-01T08:00:00+0800.
	if regDate == 0 {
		return ""
	}
	return time.Unix(regDate/1000, 0).Format("2006-01-02T15:04:05+0800")
}
