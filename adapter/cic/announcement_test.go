package cic

import (
	"github.com/tsinghua-io/api-server/resource"
	"net/http"
	"testing"
)

func TestAnnouncements(t *testing.T) {
	actual, status := adapter.Announcements("2014-2015-1-20750021-97")
	if status != http.StatusOK {
		t.Errorf("Unable to get announcements: %d", status)
		return
	}

	// Check fetched data.
	expected := &[]*resource.Announcement{
		{
			Id:        "1414652412222",
			CourseId:  "2014-2015-1-20750021-97",
			Owner:     resource.User{Name: "王媛"},
			CreatedAt: "2014-10-30",
			Priority:  0,
			Read:      true,
			Title:     "课程检索报告的要求",
			Body:      "<p><span style=\"word-wrap: break-word;\">？<span style=\"word-wrap: break-word; font-size: 7pt; font-family: &#39;Times New Roman&#39;;\">&nbsp; &nbsp; &nbsp; &nbsp;每位同学需要提交一份</span></span><span style=\"font-size: 10px;\">本课程综合检索报告。</span></p><p><span style=\"word-wrap: break-word;\">？<span style=\"word-wrap: break-word; font-size: 7pt; font-family: &#39;Times New Roman&#39;;\">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span></span>综合报告总分值为50分。<br/></p><p><span style=\"word-wrap: break-word;\">？<span style=\"word-wrap: break-word; font-size: 7pt; font-family: &#39;Times New Roman&#39;;\">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span></span>综合报告包括：选题背景分析（为什么选这个题目）、选题预期所要达成的目标、针对选题的检索步骤的细分、关于细分后的选题检索词的讨论、针对不同的选题分支的检索过程、检索效果的分析（是否符合自己的需求）以及通过对检索结果的阅读分析撰写综述，分析本课题未来可能有的研究方向等。</p><p><span style=\"word-wrap: break-word;\">？<span style=\"word-wrap: break-word; font-size: 7pt; font-family: &#39;Times New Roman&#39;;\">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span></span>综合报告是在综合利用本课所讲资源的基础上完成，一般应包括图书、期刊、学位论文、事实数据、专利文献、专科词典、百科全书等多种文献类型。</p><p><span style=\"word-wrap: break-word;\">？<span style=\"word-wrap: break-word; font-size: 7pt; font-family: &#39;Times New Roman&#39;;\">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span></span>本课旨在提高学生信息素养，鼓励大家在综合报告的最后一部分分享自己的信息使用心得，尤其是在学习、研究、社会交友等方面的网络工具，亦包括个人学习本课的体会和意见建议。</p><p><span style=\"word-wrap: break-word;\">？<span style=\"word-wrap: break-word; font-size: 7pt; font-family: &#39;Times New Roman&#39;;\">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span></span>附加题：学有余力的同学，除以上综合大报告和个人小报告之外，还可以再做一个“研究者的检索影子报告”为附加作业，附加题的分值为5分。</p><p><span style=\"word-wrap: break-word;\">？<span style=\"word-wrap: break-word; font-size: 7pt; font-family: &#39;Times New Roman&#39;;\">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span></span>研究者的检索影子报告，是像影子一样观察、复制研究者的研究行为。如，去实验室找一个学长或者一个老师，观察他们的文献检索行为，问清楚他们要检索什么，去哪里检索，如何检索的，用的什么关键词，如何处理检索结果……观察者应该像影子一样，被观察者每做一步，都应该用文字或者图像记录下来。最终形成一份“研究者的影子”报告，观察者最后应该总结：被观察者的哪些检索行为是高效的，哪些行为是与文献检索课上讲的不一致的，甚至可以给被观察者一些检索建议。<br style=\"word-wrap: break-word;\"/><br style=\"word-wrap: break-word;\"/>所有作业不限字数，但内容太单薄肯定不合适，须言之有物、条分缕析、实事求是；切勿抄袭。<br style=\"word-wrap: break-word;\"/><br style=\"word-wrap: break-word;\"/>辛苦大家，有问题随时给我邮件。</p><p><br/></p>",
		},
		{
			Id:        "1413708226270",
			CourseId:  "2014-2015-1-20750021-97",
			Owner:     resource.User{Name: "王媛"},
			CreatedAt: "2014-10-19",
			Priority:  1,
			Read:      true,
			Title:     "第4次课预习内容",
			Body:      "<p>第4次课继续讲授《文献调研》的方法，学习“文摘索引数据库”，请大家预习Web of Science平台上的web of Science核心合集的使用方法。以第一次作业中你选取的题目为例，在wos核心合集中检索相关的论文。</p><p><br/></p><p>ps：没有在网络学堂中留下邮箱地址的同学请补充一下。</p>",
		},
		{
			Id:        "1413257641943",
			CourseId:  "2014-2015-1-20750021-97",
			Owner:     resource.User{Name: "王媛"},
			CreatedAt: "2014-10-14",
			Priority:  1,
			Read:      true,
			Title:     "第三次课预习题目",
			Body:      "<p>请使用清华大学图书馆数据库导航（http://nav.lib.tsinghua.edu.cn/xport/dbdh.htm），列举你的学科最常用的数据库，并尝试使用他们</p>",
		},
		{
			Id:        "1411868112836",
			CourseId:  "2014-2015-1-20750021-97",
			Owner:     resource.User{Name: "王媛"},
			CreatedAt: "2014-09-28",
			Priority:  1,
			Read:      true,
			Title:     "第二节预习内容",
			Body:      "<p>&nbsp;请尝试获得以下两篇期刊论文的全文。</p><p><br/></p><p>1.</p><p><a href=\"http://www.ncbi.nlm.nih.gov/pubmed/21941450\" target=\"_blank\"><span style=\"color: rgb(0, 0, 255); font-family: 宋体;\">Possibility of enhanced risk of retinal neovascularization in repeated blood donors: blood donation and retinal alteration.</span></a></p><p><span style=\"font-family: 宋体;\">Rastmanesh R.</span></p><p><span style=\"font-family: 宋体;\">Int J Gen Med. 2011;4:647-56. doi: 10.2147/IJGM.S23206. Epub 2011 Sep 6.</span></p><p></p><p>2.</p><p><a href=\"http://www.ncbi.nlm.nih.gov/pubmed/19373874\" target=\"_blank\"><span style=\"color: rgb(0, 0, 255); font-family: 宋体;\">Emerging drugs: mechanism of action, mass spectrometry and <strong>doping</strong> control analysis.</span></a></p><p><span style=\"font-family: 宋体;\">Thevis M, Thomas A, Kohler M, Beuck S, Sch？nzer W.</span></p><p><span style=\"font-family: 宋体;\">J Mass Spectrom. 2009 Apr;44(4):442-60. doi: 10.1002/jms.1584.</span></p><p></p>",
		},
		{
			Id:        "1411378457399",
			CourseId:  "2014-2015-1-20750021-97",
			Owner:     resource.User{Name: "王媛"},
			CreatedAt: "2014-09-22",
			Priority:  1,
			Read:      true,
			Title:     "第一讲---预习内容",
			Body:      "<p>本课第一讲将主要帮助大家正确认识文献信息源。请大家预习以下内容。</p><p><br/></p><p>查清华大学图书馆是否有“凌晓峰.学术研究：你的成功之道.北京 : 清华大学出版社, 2012”一书。</p><p>如有，告知馆藏地、索书号和馆藏状态。</p><p>本馆是否该书的英文版本？如何使用该书电子版？</p><p><br/></p><p>请大家使用图书馆馆藏目录查阅以上信息。课堂上我会随机点名抽查预习的效果哦。</p><p><br/></p><p>明天晚上见！</p><p><br/></p><p><br/></p>",
		},
	}

	AssertDeepEqual(t, actual, expected)
}

func BenchmarkAnnouncements(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		info, status := adapter.Announcements("2014-2015-1-20750021-97")
		_ = info
		_ = status
	}
}
