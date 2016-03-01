package old

import (
	"encoding/json"
	"github.com/tsinghua-io/api-server/resource"
	"net/http"
	"testing"
)

const (
	password = ""
	username = "nxf12"
)

func TestLoginSuccuss(t *testing.T) {
	cookies, err := Login(username, password)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("Cookies received: ", cookies)
}

func TestLoginFail(t *testing.T) {
	_, err := Login("Invalidusername", "Invalidpassword")
	if err == nil {
		t.Error("Logged in using invalid username/password.")
		return
	}

	t.Log("Error received: ", err)
}

func TestPersonalInfo(t *testing.T) {
	cookies, err := Login(username, password)
	if err != nil {
		t.Error(err)
		return
	}

	adapter := New(cookies, "")
	user, status := adapter.PersonalInfo()
	if status != http.StatusOK {
		t.Errorf("Unable to get personal data: %s", err)
		return
	}

	// Check fetched data.
	expectedUser := resource.User{
		Id:         "2012011067",
		Name:       "宁雪妃",
		Type:       "本科生",
		Department: "",
		Gender:     "女",
		Email:      "1175267294@qq.com",
		Phone:      "13120098897",
	}
	if *user != expectedUser {
		t.Errorf("Incorrect data: expected %s, get %s", expectedUser, user)
		return
	}
}

func testEq(a, b []string) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestCourseIds(t *testing.T) {
	cookies, err := Login(username, password)
	if err != nil {
		t.Error(err)
		return
	}

	adapter := New(cookies, "")
	attendingIdList, err := adapter.courseIds(1)
	if err != nil {
		t.Errorf("Unable to get attending course id list: %s", err)
		return
	}

	expectedIdList := []string{
		"133593", "133106", "131792", "133107", "131777",
	}
	if !testEq(attendingIdList, expectedIdList) {
		t.Errorf("Incorrect data: excpected %s, get %s", expectedIdList, attendingIdList)
		return
	}

}

func TestCourseInfo(t *testing.T) {
	cookies, err := Login(username, password)
	if err != nil {
		t.Error(err)
		return
	}

	adapter := New(cookies, "")
	course, err := adapter.courseInfo("133593")
	if err != nil {
		t.Errorf("Unable to get course info: %s", err)
		return
	}

	j, _ := json.Marshal(*course)
	t.Logf("%s", j)
	t.Errorf("") // for debuging now
}

func TestAttending(t *testing.T) {
	cookies, err := Login(username, password)
	if err != nil {
		t.Error(err)
		return
	}

	adapter := New(cookies, "")
	courseList, status := adapter.Attending()
	if status != http.StatusOK {
		t.Errorf("Unable to get attending course info: %s", err)
		return
	}

	j, _ := json.Marshal(courseList)
	t.Logf("%s", j)
	t.Errorf("") // for debugging now
}

var fileinfos = []struct {
	path     string
	filename string
	size     int
}{
	{
		"/uploadFile/downloadFile_student.jsp?module_id=322&filePath=8arceFhSxZBoBwJb7082UV/mmNcbSN5xUe%2BpThzkK0IghF0tyxn1nKHr%2BweqOzjVD6CQMKx3SA0bx5oDxp0I024ASseHlIo8md5F3eHl5tc%3D&course_id=127759&file_id=1432738",
		"Talk_Through_533_I2S_540201987.rar",
		5254,
	},
	// Chinese characters in filename
	{
		"/uploadFile/downloadFile_student.jsp?module_id=322&filePath=0lX7YLaBEmv2fQWoiFktl6dYnDkJpaPNEM4NfmjKBHbxkaGTKEsDOQKu1bOOZIR/O36V/rEbgRs%3D&course_id=129497&file_id=1461692",
		"3.文件系统_371305032.pptx",
		436394,
	},
}

func TestParseFileInfo(t *testing.T) {
	cookies, err := Login(username, password)
	if err != nil {
		t.Error(err)
		return
	}

	adapter := New(cookies, "")
	for _, tc := range fileinfos {
		filename, size := adapter.parseFileInfo(tc.path)

		if filename != tc.filename {
			t.Errorf("Incorrect data: excpected %s, get %s", tc.filename, filename)
		}
		if size != tc.size {
			t.Errorf("Incorrect data: excpected %d, get %d", tc.size, size)
		}
	}
}

func BenchmarkLogin(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cookies, err := Login(username, password)
		_ = cookies
		_ = err
	}
}
