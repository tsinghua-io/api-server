package learn

import (
	"github.com/tsinghua-io/api-server/model"
	"github.com/tsinghua-io/api-server/util"
	"net/http"
	"testing"
)

func TestParseCourseURL(t *testing.T) {
	testSet := []struct {
		courseURL string
		courseId  string
	}{
		{"", ""},
		{"/MultiLanguage/lesson/student/course_locate.jsp?course_id=127761", "127761"},
		{"http://learn.cic.tsinghua.edu.cn/f/student/coursehome/2014-2015-1-20750021-97", "2014-2015-1-20750021-97"},
	}

	for _, testInput := range testSet {
		id := ParseCourseURL(testInput.courseURL)
		util.ExpectDeepEqual(t, id, testInput.courseId)
	}
}

func TestParseCourseName(t *testing.T) {
	testSet := []struct {
		rawName string
		name    string
		seq     string
		sem     string
	}{
		{"", "", "", ""},
		{"计算机网络(0)(2015-2016秋季学期)", "计算机网络", "0", "2015-2016-1"},
		{"   博弈论(0)(2015-2016春季学期)", "博弈论", "0", "2015-2016-2"},
		{"Matlab高级编程与工程应用(0)(2014-2015夏季学期)", "Matlab高级编程与工程应用", "0", "2014-2015-3"},
	}

	for _, testInput := range testSet {
		name, seq, sem := ParseCourseName(testInput.rawName)
		util.ExpectDeepEqual(t, name, testInput.name)
		util.ExpectDeepEqual(t, seq, testInput.seq)
		util.ExpectDeepEqual(t, sem, testInput.sem)
	}
}

func TestAllAttendedList(t *testing.T) {
	courses, status, err := ada.AllAttendedList()
	if err != nil {
		t.Fatalf("Failed to get all attended list: %s", err)
	}
	if len(courses) < 25 {
		t.Fatalf("All attended list length (%d) too small", len(courses))
	}

	util.ExpectStatus(t, status, http.StatusOK)

	testSet := []model.Course{
		{
			Id:         "132577",
			SemesterId: "2015-2016-2",
			Sequence:   "1",
			Name:       "三年级男生击剑",
		},
		{
			Id:         "108357",
			SemesterId: "2013-2014-2",
			Sequence:   "18",
			Name:       "中国近现代史纲要",
		},
		{
			Id:         "2014-2015-2-30230742-0",
			SemesterId: "2014-2015-2",
			Sequence:   "0",
			Name:       "概率论与随机过程 (1)",
		},
	}

	util.ExpectDeepEqual(t, courses[0], &testSet[0])
	util.ExpectDeepEqual(t, courses[len(courses)-1], &testSet[1])
	util.ExpectDeepEqual(t, courses[len(courses)-25], &testSet[2])
}

func BenchmarkAttended(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ada.AllAttendedList()
	}
}
