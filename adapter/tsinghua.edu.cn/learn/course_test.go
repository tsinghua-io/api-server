package learn

import (
	"github.com/tsinghua-io/api-server/adapter"
	"github.com/tsinghua-io/api-server/resource"
	"net/http"
	"testing"
)

func TestURL2CourseId(t *testing.T) {
	testSet := []struct {
		courseURL string
		courseId  string
	}{
		{"", ""},
		{"/MultiLanguage/lesson/student/course_locate.jsp?course_id=127761", "127761"},
		{"http://learn.cic.tsinghua.edu.cn/f/student/coursehome/2014-2015-1-20750021-97", "2014-2015-1-20750021-97"},
	}

	for _, testInput := range testSet {
		id, _ := URL2CourseId(testInput.courseURL)
		adapter.ExpectDeepEqual(t, id, testInput.courseId)
	}
}

func TestCourseName2Semester(t *testing.T) {
	testSet := []struct {
		Name     string
		Semester string
	}{
		{"", ""},
		{"计算机网络(0)(2015-2016秋季学期)", "2015-2016-1"},
		{"   博弈论(0)(2015-2016春季学期)", "2015-2016-2"},
		{"Matlab高级编程与工程应用(0)(2014-2015夏季学期)", "2014-2015-3"},
	}

	for _, testInput := range testSet {
		semester, _ := courseName2Semester(testInput.Name)
		adapter.ExpectDeepEqual(t, semester, testInput.Semester)
	}
}

func TestAttended(t *testing.T) {
	var courses []*resource.Course
	if status := ada.Attended("", nil, &courses); status != http.StatusOK {
		t.Errorf("Unable to get attended courses: %s", http.StatusText(status))
		return
	}

	testSet := []resource.Course{
		{
			Id:             "132577",
			Semester:       "2015-2016-2",
			CourseNumber:   "10721181",
			CourseSequence: "1",
		},
		{
			Id:             "108357",
			Semester:       "2013-2014-2",
			CourseNumber:   "10610193",
			CourseSequence: "18",
		},
		{
			Id:       "2014-2015-2-30230742-0",
			Semester: "2014-2015-2",
		},
	}

	adapter.ExpectDeepEqual(t, courses[0], &testSet[0])
	adapter.ExpectDeepEqual(t, courses[len(courses)-1], &testSet[1])
	adapter.ExpectDeepEqual(t, courses[len(courses)-25], &testSet[2])
}

func BenchmarkAttended(b *testing.B) {
	var courses []*resource.Course

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ada.Attended("", nil, &courses)
	}
}

func TestCourseIdMap(t *testing.T) {
	actual := make(map[string]string)
	if status := ada.CourseIdMap("", nil, actual); status != http.StatusOK {
		t.Errorf("Unable to get attended courses: %s", http.StatusText(status))
		return
	}

	adapter.ExpectDeepEqual(t, actual["2014-2015-2-40260202-0"], "")
	adapter.ExpectDeepEqual(t, actual["2014-2015-2-30230711-2"], "123510")
	adapter.ExpectDeepEqual(t, actual["2013-2014-2-10610193-18"], "108357")
}
