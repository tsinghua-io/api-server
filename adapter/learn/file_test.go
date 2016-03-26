package learn

import (
	"github.com/tsinghua-io/api-server/adapter"
	"github.com/tsinghua-io/api-server/resource"
	"net/http"
	"testing"
)

func TestCourseFiles(t *testing.T) {
	var actual []*resource.File
	if status := ada.Files("127756", nil, &actual); status != http.StatusOK {
		t.Errorf("Unable to get files: %s", http.StatusText(status))
		return
	}

	// Check fetched data.
	tab1file1 := &resource.File{
		Id:          "1426717",
		CourseId:    "127756",
		CreatedAt:   "2015-09-13",
		Title:       "第1讲 操作系统引论(1)",
		Description: "",
		Category:    []string{"电子教案"},
		Filename:    "01操作系统引论(1)_178207362.pdf",
		Size:        2481982,
		DownloadURL: "https://learn.tsinghua.edu.cn/uploadFile/downloadFile_student.jsp?module_id=322\u0026filePath=N4Tel3ukBcf0P%2BxFdYeeoHN1562AESxTOYGA60sn5xpe9dkSxsvYaLcsha/n4FAm\u0026course_id=127756\u0026file_id=1426717",
	}

	tab2file1 := &resource.File{
		Id:          "1429268",
		CourseId:    "127756",
		CreatedAt:   "2015-09-14",
		Title:       "现代操作系统（第三版）英文版",
		Description: "",
		Category:    []string{"补充资料"},
		Filename:    "modern_operating_systems_3rd_edition_tanenbaum_171607868.pdf",
		Size:        16402536,
		DownloadURL: "https://learn.tsinghua.edu.cn/uploadFile/downloadFile_student.jsp?module_id=322\u0026filePath=tAY6dAg1JH0INPji6josGqd/QGxTAbvadCvv0EfUnWw1ilm2qC/RZXchbtqC3FfuswFdhSOzbohNc8dms8TKZOiOp0KJm7vo8kXwiCOpbiSvRRDvZlFfPmDX4MQCdxbueQFju3W3qmM%3D\u0026course_id=127756\u0026file_id=1429268",
	}

	adapter.AssertDeepEqual(t, actual[0], tab1file1)
	adapter.AssertDeepEqual(t, actual[16], tab2file1)
}

func BenchmarkCourseFiles(b *testing.B) {
	var files []*resource.File

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ada.Files("127756", nil, &files)
	}
}
