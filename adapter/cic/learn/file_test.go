package learn

import (
	"github.com/tsinghua-io/api-server/model"
	"github.com/tsinghua-io/api-server/util"
	"net/http"
	"testing"
)

func TestFiles(t *testing.T) {
	actual, status, err := ada.Files("2014-2015-1-20750021-97")
	if err != nil {
		t.Fatalf("Failed to get files: %s", err)
	}

	util.ExpectStatus(t, status, http.StatusOK)

	// Check fetched data.
	expected := []*model.File{
		{
			Id:          "2004980851_2014-2015-1-20750021-97_KJ_1411486091",
			CourseId:    "2014-2015-1-20750021-97",
			Owner:       &model.User{Id: "2004980851", Name: "王媛"},
			CreatedAt:   "2014-09-23T23:28:13+0800",
			Title:       "全面认识文献信息源1",
			Description: "本讲中提到的工具：图书馆主页、馆藏目录、超星电子图书、读秀学术搜索、FirstSearch中的WorldCat联合目录。",
			Category:    []string{"课程文件", "电子教案"},
			Filename:    "文献检索与利用（理工类）-全面认识文献信息源.pptx",
			Size:        7935551,
			DownloadURL: "http://learn.cic.tsinghua.edu.cn/b/resource/downloadFileStream/2004980851_2014-2015-1-20750021-97_KJ_1411486091",
		},
		{
			Id:          "2004980851_2014-2015-1-20750021-97_KJ_1413292186",
			CourseId:    "2014-2015-1-20750021-97",
			Owner:       &model.User{Id: "2004980851", Name: "王媛"},
			CreatedAt:   "2014-10-14T21:09:47+0800",
			Title:       "全面认识文献信息源2",
			Description: "",
			Category:    []string{"课程文件", "电子教案"},
			Filename:    "文献检索与利用（理工类）-全面认识文献信息源2.pptx",
			Size:        3678507,
			DownloadURL: "http://learn.cic.tsinghua.edu.cn/b/resource/downloadFileStream/2004980851_2014-2015-1-20750021-97_KJ_1413292186",
		},
		{
			Id:          "2004980851_2014-2015-1-20750021-97_KJ_1413292258",
			CourseId:    "2014-2015-1-20750021-97",
			Owner:       &model.User{Id: "2004980851", Name: "王媛"},
			CreatedAt:   "2014-10-14T21:10:58+0800",
			Title:       "文献调研1",
			Description: "",
			Category:    []string{"课程文件", "电子教案"},
			Filename:    "文献检索与利用（3）--文献调研1.pptx",
			Size:        5154562,
			DownloadURL: "http://learn.cic.tsinghua.edu.cn/b/resource/downloadFileStream/2004980851_2014-2015-1-20750021-97_KJ_1413292258",
		},
		{
			Id:          "2004980851_2014-2015-1-20750021-97_KJ_1414651951",
			CourseId:    "2014-2015-1-20750021-97",
			Owner:       &model.User{Id: "2004980851", Name: "王媛"},
			CreatedAt:   "2014-10-30T14:52:32+0800",
			Title:       "文献调研SCI",
			Description: "",
			Category:    []string{"课程文件", "电子教案"},
			Filename:    "文献检索与利用（3）--文献调研2--SCI.pptx",
			Size:        4807443,
			DownloadURL: "http://learn.cic.tsinghua.edu.cn/b/resource/downloadFileStream/2004980851_2014-2015-1-20750021-97_KJ_1414651951",
		},
		{
			Id:          "2004980851_2014-2015-1-20750021-97_KJ_1414652071",
			CourseId:    "2014-2015-1-20750021-97",
			Owner:       &model.User{Id: "2004980851", Name: "王媛"},
			CreatedAt:   "2014-10-30T14:54:31+0800",
			Title:       "文献调研EI",
			Description: "",
			Category:    []string{"课程文件", "电子教案"},
			Filename:    "文献检索与利用（3）--文献调研3--EI.pptx",
			Size:        2275654,
			DownloadURL: "http://learn.cic.tsinghua.edu.cn/b/resource/downloadFileStream/2004980851_2014-2015-1-20750021-97_KJ_1414652071",
		},
	}

	util.ExpectDeepEqual(t, actual, expected)
}

func BenchmarkFiles(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ada.Files("2014-2015-1-20750021-97")
	}
}
