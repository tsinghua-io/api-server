package learn

import (
	"net/http"
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
