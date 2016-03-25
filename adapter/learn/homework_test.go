package old

import (
	"github.com/tsinghua-io/api-server/resource"
	"net/http"
	"testing"
)

func TestCourseHomework(t *testing.T) {
	actual, status := ada.CourseHomework("127743", nil)
	if status != http.StatusOK {
		t.Errorf("Unable to get homework: %s", http.StatusText(status))
		return
	}

	// Check fetched data.
	expected := []*resource.Homework{
		{
			Id:       "663857",
			CourseId: "127743",
			BeginAt:  "2015-11-02",
			DueAt:    "2015-11-05",
			Title:    "第六次作业",
			Body:     "    第六次作业有部分内容需要Matlab仿真，所以如果有同学想提交电子版作业，请把非仿真的部分也拍照一并提交上来，这样就可以不必提交纸版作业了。\n    如果有同学要提交纸板作业，请把仿真内容打印出来，在课堂上交给老师！\n",
			Submissions: []*resource.Submission{
				{
					HomeworkId: "663857",
					Late:       false,
					Body:       "",
					Attachment: &resource.Attachment{
						Filename:    "2013011187_663857_873609575_p4.pdf",
						Size:        2820876,
						DownloadUrl: "https://learn.tsinghua.edu.cn/uploadFile/downloadFile.jsp?module_id=322\u0026course_id=127743\u0026filePath=Ui6dWfN3E23iy92Lm3GqLolVIj%2Bu5tfsytg7jRwlhOqxaULAEWF80pMjNsAbeGoNLbxf932lsPZSeaPFGySxlzqYaxPQuvWF9JTL%2B1WuOg4%3D",
					},
					MarkedAt: "2015-11-06",
					Mark:     NewFloat32(9.5),
					Comment:  "第二题分析和结论正确，公式有问题。",
				},
			},
		},
		{
			Id:       "667021",
			CourseId: "127743",
			BeginAt:  "2015-11-15",
			DueAt:    "2015-11-30",
			Title:    "大作业",
			Body:     "作业说明件附件，参考论文和数据请从“课程文件”中下载。\n\n",
			Attachment: &resource.Attachment{
				Filename:    "625602385_2_2015年《数字信号处理》课程大作业.pdf",
				Size:        112055,
				DownloadUrl: "https://learn.tsinghua.edu.cn/uploadFile/downloadFile.jsp?module_id=322\u0026course_id=127743\u0026filePath=7D5eM/3uxuWgUscnZFe5xYFRwCtzmT3Nd4b8XfYdVt9QXP6jW0X3Mw6gr2ogb0t8bD67/q7AeDDvr3x32279mpdW6Tj5nS6ysO1fFyPcUzk%3D",
			},
			Submissions: []*resource.Submission{
				{
					HomeworkId: "667021",
					Late:       false,
					Body:       "",
					Attachment: &resource.Attachment{
						Filename:    "2013011187_667021_531504538_report.pdf",
						Size:        98640,
						DownloadUrl: "https://learn.tsinghua.edu.cn/uploadFile/downloadFile.jsp?module_id=322\u0026course_id=127743\u0026filePath=Ui6dWfN3E23iy92Lm3GqLolVIj%2Bu5tfsVetxSI%2BmeI5zL/GWM0GkxzPppRm00efUNVY7MLZOt3A1jm56tM3YdeAlZMTa30DiABpxaPmB1YI%3D",
					},
				},
			},
		},
		{
			Id:       "669225",
			CourseId: "127743",
			BeginAt:  "2015-11-24",
			DueAt:    "2015-11-26",
			Title:    "第九次作业",
			Body:     "",
		},
		{
			Id:       "670485",
			CourseId: "127743",
			BeginAt:  "2015-11-30",
			DueAt:    "2015-12-04",
			Title:    "第十次作业",
			Body:     "请未交作业的同学尽快提交作业！",
		},
		{
			Id:       "672559",
			CourseId: "127743",
			BeginAt:  "2015-12-09",
			DueAt:    "2015-12-11",
			Title:    "第十一次作业",
			Body:     "",
		},
		{
			Id:       "674369",
			CourseId: "127743",
			BeginAt:  "2015-12-15",
			DueAt:    "2015-12-18",
			Title:    "第十二次作业",
			Body:     "",
		},
		{
			Id:       "675999",
			CourseId: "127743",
			BeginAt:  "2015-12-22",
			DueAt:    "2015-12-25",
			Title:    "第十三次作业",
			Body:     "",
		},
		{
			Id:       "677925",
			CourseId: "127743",
			BeginAt:  "2015-12-30",
			DueAt:    "2016-01-01",
			Title:    "第十四次作业",
			Body:     "",
			Submissions: []*resource.Submission{
				{
					HomeworkId: "677925",
					Late:       false,
					Body:       "助教你好，这是最后的三次作业（对于课件12，13，14，15），一起交上来了！",
					Attachment: &resource.Attachment{
						Filename:    "2013011187_677925_351002502_课件12-15对应的作业.zip",
						Size:        13601942,
						DownloadUrl: "https://learn.tsinghua.edu.cn/uploadFile/downloadFile.jsp?module_id=322\u0026course_id=127743\u0026filePath=Ui6dWfN3E23iy92Lm3GqLolVIj%2Bu5tfsAqZGEzdQ71%2BH77AyehAScSLl2e2n5PeFBbCCeknxk/vScQfKcoZU4o%2B/fgVxuHWOFnv5a2X8eUU%3D",
					},
					MarkedAt: "2016-01-03",
					Mark:     NewFloat32(9.5),
					Comment:  "前面三次各8分",
				},
			},
		},
	}

	AssertDeepEqual(t, actual, expected)
}

func BenchmarkCourseHomework(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		info, status := ada.CourseHomework("127743", nil)
		_ = info
		_ = status
	}
}
