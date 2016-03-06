// Adapters

package adapter

import (
	"github.com/golang/glog"
	"github.com/tsinghua-io/api-server/resource"
	"io"
	"net/http"
	"time"
)

type Adapter interface {
	Profile(username string, params map[string]string) (user *resource.User, status int)
	Attended(username string, params map[string]string) (courses []*resource.Course, status int)

	CourseAnnouncements(courseId string, params map[string]string) (announcements []*resource.Announcement, status int)
	CourseFiles(courseId string, params map[string]string) (files []*resource.File, status int)
	CourseHomework(courseId string, params map[string]string) (homework []*resource.Homework, status int)
}

type Parser interface {
	Parse(reader io.Reader, info interface{}) error
}

// FetchInfo fetches info from url using HTTP GET/POST, and then parse it
// using the given parser. A HTTP status code is returned to indicate the
// result.
func FetchInfo(client *http.Client, url string, method string, p Parser, info interface{}) (status int) {
	// Fetch data from url.
	glog.Infof("Fetching data from %s", url)

	var resp *http.Response
	var err error
	t_send := time.Now()

	switch method {
	case "GET":
		resp, err = client.Get(url)
	case "POST":
		resp, err = client.Post(url, "application/x-www-form-urlencoded", nil)
	default:
		glog.Errorf("Unknown method to fetch info: %s", method)
		return http.StatusInternalServerError
	}
	if err != nil {
		glog.Errorf("Unable to fetch info from %s: %s", url, err)
		return http.StatusBadGateway
	}
	defer resp.Body.Close()

	t_receive := time.Now()
	glog.Infof("Fetched data from %s (%s)", url, t_receive.Sub(t_send))

	// Parse the data.
	if err := p.Parse(resp.Body, info); err != nil {
		glog.Errorf("Unable to parse data received from %s: %s", url, err)
		return http.StatusInternalServerError
	}

	glog.Infof("Parsed data from %s (%s)", url, time.Since(t_receive))

	// We are safe.
	return http.StatusOK
}
