package learn

import (
	"github.com/tsinghua-io/api-server/adapter"
	"github.com/tsinghua-io/api-server/resource"
	"net/http"
	"testing"
)

func TestCourseAnnouncements(t *testing.T) {
	var actual []*resource.Announcement
	if status := ada.Announcements("103048", nil, &actual); status != http.StatusOK {
		t.Errorf("Unable to get announcements: %s", http.StatusText(status))
		return
	}

	// Check fetched data.
	expected := []*resource.Announcement{
		{
			Id:        "1535567",
			CourseId:  "103048",
			Owner:     &resource.User{Name: "孙甲松老师"},
			CreatedAt: "2013-12-31",
			Priority:  1,
			Title:     "期末考试时间地点",
			Body:      " <font size=\"6\">考试时间：1月18日上午8:00~9:00<br/>考试地点：一教101     学号：2012011072~2013011167<br/>                  一教104     学号：2013011168~2013080069</font> \n\t\t\t",
		},
		{
			Id:        "1481094",
			CourseId:  "103048",
			Owner:     &resource.User{Name: "陈权崎"},
			CreatedAt: "2013-10-16",
			Priority:  0,
			Title:     "捡到一个U盘",
			Body:      " 黄永峰老师的程序设计助教陶怀舟在机房捡到一个U盘，有丢失U盘的同学可以去罗姆楼8-203或者周五晚上上机时间到机房找陶怀舟查看。  \n\t\t\t",
		},
		{
			Id:        "1466904",
			CourseId:  "103048",
			Owner:     &resource.User{Name: "孙甲松老师"},
			CreatedAt: "2013-09-26",
			Priority:  1,
			Title:     "《C程序设计教程》马上到教材中心",
			Body:      " <font size=\"5\">教材中心通知：《C程序设计教程》今天下午4点后到教材中心，大家可以以班级或者个人去购买。<br/>请相互转告！</font> \n\t\t\t",
		},
		{
			Id:        "1466172",
			CourseId:  "103048",
			Owner:     &resource.User{Name: "孙甲松老师"},
			CreatedAt: "2013-09-25",
			Priority:  1,
			Title:     "注意：C语言02课件已经更新",
			Body:      " <p><font size=\"4\">注意：C语言02课件已经更新，第2次实验中，98.6875改为97.6875，以与教材和ppt中所给示例的值保持一致。抱歉给大家添乱！<br/><br/>1.编写程序打印查看<font color=\"#ff0000\">97.6875</font>的double型值在计算机内的IEEE 754 存储格式按字节的十六进制值，验证是否与教材所给结果一致。同时也打印查看97.6875的float型值在计算机内存储按字节的十六进制值。</font></p> \n\t\t\t",
		},
		{
			Id:        "1463882",
			CourseId:  "103048",
			Owner:     &resource.User{Name: "孙甲松老师"},
			CreatedAt: "2013-09-23",
			Priority:  1,
			Title:     "分组与助教联系方式",
			Body:      " <font size=\"4\">所有同学被分成4组，每人有一个组号，由学号所在的范围决定，与班号无关。<br/>第1组：2012011072~2013011110<font size=\"4\"><br/>第2组：2013011112~2013011167<br/></font><font size=\"4\">第3组：2013011168~2013011234</font><font size=\"4\"><br/>第4组：2013011236~2013080069<br/></font><br/>以下4位助教分别负责批改第1~第4组作业：<br/>第1组（2012011072~2013011110）：陈权崎  13810323433    </font><a href=\"mailto:quanqi.chen@gmail.com\"><font size=\"4\">quanqi.chen@gmail.com</font></a><font size=\"4\">     <br/>第2组（2013011112~2013011167）：李雪      13488687568    </font><a style=\"FONT-SIZE: 14px; FONT-FAMILY: arial,sans-serif\" target=\"_blank\" href=\"mailto:xue-li11@mails.tsinghua.edu.cn\"><font size=\"4\">xue-li11@mails.tsinghua.edu.cn</font></a><font size=\"4\">     <br/>第3组（2013011168~2013011234）：张蕤      15120083955    </font><a href=\"mailto:zr1174@163.com\"><font size=\"4\">zr1174@163.com</font></a><font size=\"4\">     <br/>第4组（2013011236~2013080069）：王斌      15120001254    </font><a href=\"mailto:wb.th08@gmail.com\"><font size=\"4\">wb.th08@gmail.com</font></a><font size=\"4\">  <br/><br/>有问题请与他们联系。</font><br/> \n\t\t\t",
		},
		{
			Id:        "1463860",
			CourseId:  "103048",
			Owner:     &resource.User{Name: "孙甲松老师"},
			CreatedAt: "2013-09-23",
			Priority:  1,
			Title:     "提交作业或者实验报告后请检查附件大小！",
			Body:      " <font size=\"4\">各位同学，提交作业或者实验报告后请检查附件大小！如果出现文件长度为0或者“无附件”，说明提交不成功，请重新提交! 切记！</font> \n\t\t\t",
		},
		{
			Id:        "1460655",
			CourseId:  "103048",
			Owner:     &resource.User{Name: "孙甲松老师"},
			CreatedAt: "2013-09-18",
			Priority:  1,
			Title:     "关于前两次纸值作业",
			Body:      " <font size=\"4\">第一、二次纸值作业，最好以班级为单位收齐，下次上课时由课代表交给我，我转交给助教。<br/>其余的作业和实验，请在网络学堂中的“课程作业”里提交电子版。请相互转告！</font> \n\t\t\t",
		},
		{
			Id:        "1460638",
			CourseId:  "103048",
			Owner:     &resource.User{Name: "孙甲松老师"},
			CreatedAt: "2013-09-18",
			Priority:  1,
			Title:     "实验写报告提交到网络学堂中",
			Body:      " <font size=\"3\">各位同学，<br/><br/>以后的实验以及第三次以后的作业，都写一个电子文档 (.doc 或.docx) <strong>提交到网络学堂中</strong>，</font><font size=\"3\"><strong>注意：不要发到我的邮箱中！<br/></strong>因为助教负责批改作业，不提交到网络学堂中他们看不到，无法给成绩。<br/>即使已经用邮件发给我的同学，也请到网络学堂中提交。请相互转告！<br/><br/>文档内容：<br/>6.     源程序： 。。。。。。。。<br/>       运行结果：（截图或者结果）<br/>7.     源程序： 。。。。。。。。<br/>       运行结果：（截图或者结果）<br/><br/>如果是实验，要求写<strong>实验报告</strong>，请按<strong>实验报告的格式</strong>写：<br/>题目：。。。。。。<br/>源程序：。。。。。。。。<br/>运行结果：（截图或者结果）<br/>结果分析与总结：.。。。。。。。。。</font> \n\t\t\t",
		},
	}

	adapter.AssertDeepEqual(t, actual, expected)
}

func BenchmarkCourseAnnouncements(b *testing.B) {
	var announcements []*resource.Announcement

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ada.Announcements("103048", nil, &announcements)
	}
}
