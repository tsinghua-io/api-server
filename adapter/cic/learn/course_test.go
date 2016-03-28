package learn

import (
	"github.com/tsinghua-io/api-server/model"
	"github.com/tsinghua-io/api-server/util"
	"net/http"
	"testing"
)

func TestAllAttended(t *testing.T) {
	if courses, status := ada.AllAttended(false); status != http.StatusOK {
		t.Fatalf("Unable to get attended courses: %s", http.StatusText(status))
	}

	// Just test the last course.
	if len(courses) == 0 {
		t.Fatalf("Got no courses.")
	}
	actual := courses[len(courses)-1]
	expected := &model.Course{
		Id:             "2013-2014-1-00640252-96",
		Semester:       "2013-2014-1",
		CourseNumber:   "00640252",
		CourseSequence: "96",
		Name:           "英语报刊选读",
		Credit:         2,
		Hour:           32,
		Description:    "本课程将引导学生阅读当代英美报刊不同体裁的文章。使学生初步了解英美报刊文章的特点，学会识别不同体裁。在阅读能力提高的基础上，增强用英语表述自我观点的能力，从而加强批判性思维的能力。所选主题包括经济、环境、战争、科技、教育、社会、政府和体育，以及各类时事。要求学生自选5篇社论/专栏等文章写读书报告，完成5篇新闻总结。",
		TimeLocations: []*model.TimeLocation{
			{
				Weeks:       "全周",
				DayOfWeek:   4,
				PeriodOfDay: 1,
				Location:    "六教6B105",
			},
		},
		Teachers: []*model.User{
			{
				Id:         "L064533",
				Name:       "Andrew Backe",
				Department: "外国语言文学系",
				Gender:     "男",
			},
		},
	}

	util.AssertDeepEqual(t, actual, expected)
}

func BenchmarkAllAttended(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ada.AllAttended(false)
	}
}
