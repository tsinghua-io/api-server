package learn

import (
	"github.com/tsinghua-io/api-server/model"
	"github.com/tsinghua-io/api-server/util"
	"net/http"
	"testing"
)

func TestAssignments(t *testing.T) {
	actual, status, err := ada.Assignments("2014-2015-1-20750021-97")
	if err != nil {
		t.Fatalf("Failed to get homeworks: %s", err)
	}

	util.ExpectStatus(t, status, http.StatusOK)

	// Check fetched data.
	expected := []*model.Assignment{
		{
			Id:          "58093",
			CourseId:    "2014-2015-1-20750021-97",
			CreatedAt:   "",
			BeginAt:     "2014-09-23T10:02:50+0800",
			DueAt:       "2014-09-23T18:59:59+0800",
			Title:       "第一节课预习作业：已经发了课程公告，为防止大家没看到",
			Body:        "<p style=\"margin-top: 0px; margin-bottom: 0px; padding: 0px; font-family: Tahoma, Helvetica, Arial, 微软雅黑, sans-serif; font-size: 12px; line-height: 22px; white-space: normal; background-color: rgb(255, 255, 255);\">重要：预习作业每人都要做，但不用交，我会课堂上随机抽查预习效果。</p><p style=\"margin-top: 0px; margin-bottom: 0px; padding: 0px; font-family: Tahoma, Helvetica, Arial, 微软雅黑, sans-serif; font-size: 12px; line-height: 22px; white-space: normal; background-color: rgb(255, 255, 255);\"><br/></p><p style=\"margin-top: 0px; margin-bottom: 0px; padding: 0px; font-family: Tahoma, Helvetica, Arial, 微软雅黑, sans-serif; font-size: 12px; line-height: 22px; white-space: normal; background-color: rgb(255, 255, 255);\">本课第一讲将主要帮助大家正确认识文献信息源。请大家预习以下内容。</p><p style=\"margin-top: 0px; margin-bottom: 0px; padding: 0px; font-family: Tahoma, Helvetica, Arial, 微软雅黑, sans-serif; font-size: 12px; line-height: 22px; white-space: normal; background-color: rgb(255, 255, 255);\"><br/></p><p style=\"margin-top: 0px; margin-bottom: 0px; padding: 0px; font-family: Tahoma, Helvetica, Arial, 微软雅黑, sans-serif; font-size: 12px; line-height: 22px; white-space: normal; background-color: rgb(255, 255, 255);\">查清华大学图书馆是否有“凌晓峰.学术研究：你的成功之道.北京 : 清华大学出版社, 2012”一书。</p><p style=\"margin-top: 0px; margin-bottom: 0px; padding: 0px; font-family: Tahoma, Helvetica, Arial, 微软雅黑, sans-serif; font-size: 12px; line-height: 22px; white-space: normal; background-color: rgb(255, 255, 255);\">如有，告知馆藏地、索书号和馆藏状态。</p><p style=\"margin-top: 0px; margin-bottom: 0px; padding: 0px; font-family: Tahoma, Helvetica, Arial, 微软雅黑, sans-serif; font-size: 12px; line-height: 22px; white-space: normal; background-color: rgb(255, 255, 255);\">本馆是否该书的英文版本？如何使用该书电子版？</p><p style=\"margin-top: 0px; margin-bottom: 0px; padding: 0px; font-family: Tahoma, Helvetica, Arial, 微软雅黑, sans-serif; font-size: 12px; line-height: 22px; white-space: normal; background-color: rgb(255, 255, 255);\"><br/></p><p style=\"margin-top: 0px; margin-bottom: 0px; padding: 0px; font-family: Tahoma, Helvetica, Arial, 微软雅黑, sans-serif; font-size: 12px; line-height: 22px; white-space: normal; background-color: rgb(255, 255, 255);\">请大家使用图书馆馆藏目录查阅以上信息。课堂上我会随机点名抽查预习的效果哦。</p><p style=\"margin-top: 0px; margin-bottom: 0px; padding: 0px; font-family: Tahoma, Helvetica, Arial, 微软雅黑, sans-serif; font-size: 12px; line-height: 22px; white-space: normal; background-color: rgb(255, 255, 255);\"><br/></p><p style=\"margin-top: 0px; margin-bottom: 0px; padding: 0px; font-family: Tahoma, Helvetica, Arial, 微软雅黑, sans-serif; font-size: 12px; line-height: 22px; white-space: normal; background-color: rgb(255, 255, 255);\">晚上见！</p><p><br/></p>",
			Attachment:  nil,
			Submissions: nil,
		},
		{
			Id:         "71056",
			CourseId:   "2014-2015-1-20750021-97",
			CreatedAt:  "2014-10-14T21:18:04+0800",
			BeginAt:    "2014-10-14T00:00:00+0800",
			DueAt:      "2014-10-20T23:59:59+0800",
			Title:      "第一次作业：电子图书与检索式编写",
			Body:       "<p>请于10月20日前提交作业，电子版即可，可将两个题目做在一个ppt或word中，推荐使用ppt呈现。有不明白的请给我或助教发邮件。<br/></p>",
			Attachment: nil,
			Submissions: []*model.Submission{
				{
					Owner:        &model.User{Id: "2013011187"},
					AssignmentId: "71056",
					CreatedAt:    "2014-10-20T23:50:27+0800",
					Late:         false,
					Body:         "",
					Attachment: &model.Attachment{
						Filename:    "作业1.pptx",
						Size:        315198,
						DownloadURL: "http://learn.cic.tsinghua.edu.cn/b/resource/downloadFileStream/2013011187_2014-2015-1-20750021-97_ZY_1413820224",
					},
					MarkedBy:          &model.User{Name: "汤娇"},
					MarkedAt:          "2014-10-23T23:12:39+0800",
					Mark:              util.NewFloat32(11),
					Comment:           "<p>作业一,缺少分享一部分,此次不扣分，请再试一下如何将下载的图书拷贝给别人使用,另外关于超星电子图书阅读器的注册和下载过程的展示过于简略,;作业二,中文检索式连接符的使用正确,但是检索词的选择过于生僻和口语化,导致你的检索结果没有,英文检索式错误,请注意括号需要成对出现和检索词的选择.</p>",
					CommentAttachment: nil,
				},
			},
		},
		{
			Id:        "76121",
			CourseId:  "2014-2015-1-20750021-97",
			CreatedAt: "2014-10-22T17:19:10+0800",
			BeginAt:   "2014-10-22T00:00:00+0800",
			DueAt:     "2014-11-01T23:59:59+0800",
			Title:     "第二次作业——SCI",
			Body:      "<p>课件我传不上去，只要放在印象笔记中与大家共享，请大家通过这个链接<a href=\"http://app.yinxiang.com/l/AAmpP7Z7r6lKnqrhVrtaVajvcJNd__J4xEA/\">http://app.yinxiang.com/l/AAmpP7Z7r6lKnqrhVrtaVajvcJNd__J4xEA/</a>&nbsp;来下载，以前和以后的课件我多会放在这里。请大家按时完成作业。</p>",
			Attachment: &model.Attachment{
				Filename:    "文献检索与利用-第二次作业——SCI.docx",
				Size:        15832,
				DownloadURL: "http://learn.cic.tsinghua.edu.cn/b/resource/downloadFileStream/2004980851_2014-2015-1-20750021-97_ZY_1413969147",
			},
			Submissions: []*model.Submission{
				{
					Owner:        &model.User{Id: "2013011187"},
					AssignmentId: "76121",
					CreatedAt:    "2014-11-01T23:58:44+0800",
					Late:         false,
					Body:         "",
					Attachment: &model.Attachment{
						Filename:    "文献检索与利用-第二次作业——SCI.docx",
						Size:        250204,
						DownloadURL: "http://learn.cic.tsinghua.edu.cn/b/resource/downloadFileStream/2013011187_2014-2015-1-20750021-97_ZY_1414857519",
					},
					MarkedBy:          &model.User{Name: "汤娇"},
					MarkedAt:          "2014-11-11T18:26:26+0800",
					Mark:              util.NewFloat32(8),
					Comment:           "<p><span style=\"font-size:14px;font-family:宋体\">李思涵,你好:</span></p><p><span style=\"font-size:14px;font-family:宋体\"><br/></span></p><p><span style=\"font-size:14px;font-family:宋体\">&nbsp; &nbsp;第二次作业已阅,完成得不是很理想,你可以通过优化一下你的检索式,来获得最适合你的检索结果,可以适当放开你的检索条件;否则你第三题你如何完成这部分的作业;有问题的话,请联系王老师和我.</span></p><p><span style=\"font-size:14px;font-family:宋体\"><br/></span></p><p><span style=\"font-size:14px;font-family:宋体\">&nbsp; 继续努力,将后续的作业做好!</span></p><p><span style=\"font-size:14px;font-family:宋体\"><br/></span></p><p><span style=\"font-size:14px;font-family:宋体\">汤娇</span></p>",
					CommentAttachment: nil,
				},
			},
		},
		{
			Id:        "88053",
			CourseId:  "2014-2015-1-20750021-97",
			CreatedAt: "2014-10-30T14:56:19+0800",
			BeginAt:   "2014-10-30T00:00:00+0800",
			DueAt:     "2014-11-14T23:59:59+0800",
			Title:     "EI",
			Body:      "",
			Attachment: &model.Attachment{
				Filename:    "EI作业.doc",
				Size:        32256,
				DownloadURL: "http://learn.cic.tsinghua.edu.cn/b/resource/downloadFileStream/2004980851_2014-2015-1-20750021-97_ZY_1414652158",
			},
			Submissions: []*model.Submission{
				{
					Owner:        &model.User{Id: "2013011187"},
					AssignmentId: "88053",
					CreatedAt:    "2014-11-14T11:17:18+0800",
					Late:         false,
					Body:         "",
					Attachment: &model.Attachment{
						Filename:    "EI作业.doc",
						Size:        36352,
						DownloadURL: "http://learn.cic.tsinghua.edu.cn/b/resource/downloadFileStream/2013011187_2014-2015-1-20750021-97_ZY_1415935032",
					},
					MarkedBy:          &model.User{Name: "汤娇"},
					MarkedAt:          "2014-11-19T21:11:44+0800",
					Mark:              util.NewFloat32(14),
					Comment:           "<p style=\"TEXT-ALIGN: left; MARGIN: 0px 0px 5px\"><span style=\"FONT-FAMILY: 宋体; COLOR: black; FONT-SIZE: 16px\">李思涵，你好：</span></p><p style=\"TEXT-ALIGN: left; MARGIN: 5px 0px\"><span style=\"FONT-FAMILY: 宋体; COLOR: black; FONT-SIZE: 16px\">&nbsp;&nbsp;&nbsp;&nbsp; </span></p><p style=\"TEXT-ALIGN: left; MARGIN: 5px 0px\"><span style=\"FONT-FAMILY: 宋体; COLOR: black; FONT-SIZE: 16px\">&nbsp;&nbsp;&nbsp;&nbsp; EI</span><span style=\"FONT-FAMILY: 宋体; COLOR: black; FONT-SIZE: 16px\">作业已阅，总体完成得不错，过程和思路也比较清晰，但是按照你的expert search，应该只有4个检索结果。</span></p><p style=\"TEXT-ALIGN: left; MARGIN: 5px 0px\"><span style=\"FONT-FAMILY: 宋体; COLOR: black; FONT-SIZE: 16px\">&nbsp;&nbsp;&nbsp;&nbsp; </span><span style=\"FONT-FAMILY: 宋体; COLOR: black; FONT-SIZE: 16px\">继续努力，将检索报告好好完成！</span><span style=\"FONT-FAMILY: &#39;Arial&#39;,&#39;sans-serif&#39;; FONT-SIZE: 16px\">&nbsp;</span></p><p style=\"TEXT-ALIGN: left; MARGIN: 5px 0px\"><span style=\"FONT-FAMILY: 宋体; COLOR: black; FONT-SIZE: 16px\">&nbsp;</span><span style=\"FONT-FAMILY: 宋体; COLOR: black; FONT-SIZE: 16px\">汤娇</span></p><p></p>",
					CommentAttachment: nil,
				},
			},
		},
		{
			Id:         "92095",
			CourseId:   "2014-2015-1-20750021-97",
			CreatedAt:  "2014-11-15T15:27:16+0800",
			BeginAt:    "2014-11-15T00:00:00+0800",
			DueAt:      "2014-12-01T23:59:59+0800",
			Title:      "综合检索报告提交专用",
			Body:       "",
			Attachment: nil,
			Submissions: []*model.Submission{
				{
					Owner:        &model.User{Id: "2013011187"},
					AssignmentId: "92095",
					CreatedAt:    "2014-12-02T13:21:56+0800",
					Late:         true,
					Body:         "",
					Attachment: &model.Attachment{
						Filename:    "final.pdf",
						Size:        223372,
						DownloadURL: "http://learn.cic.tsinghua.edu.cn/b/resource/downloadFileStream/2013011187_2014-2015-1-20750021-97_ZY_1417497712",
					},
					MarkedBy:          &model.User{Name: "王媛"},
					MarkedAt:          "2015-02-01T22:16:25+0800",
					Mark:              util.NewFloat32(38),
					Comment:           "",
					CommentAttachment: nil,
				},
			},
		},
	}

	util.ExpectDeepEqual(t, actual, expected)
}

func BenchmarkAssignments(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ada.Assignments("2014-2015-1-20750021-97")
	}
}
