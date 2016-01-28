//Definitions of resources.

package resources

type SUser struct {
	Id         string
	Name       string
	User_type  string
	Department string
}

type User struct {
	SUser
	Gender string
	Email  string
	Phone  string
	Class  string
}

type SCourse struct {
	Id          string
	Name        string
	Teacher     SUser
	Coteacher   []SUser
	School_year string
	Semester    string
}

type Course struct {
	SCourse
	Course_number   string
	Course_sequence string
	Credit          int
	Hour            int
	Description     string
	Students_count  int
}

type SAnnounce struct {
	Course_id  string
	Number     int
	Title      string
	Owner      SUser
	Created_at string
	Updated_at string
	Important  bool
}

type Announce struct {
	SAnnounce
	Body string
}

type File struct {
	Course_id    string
	Number       int
	Category     []string
	Title        string
	Description  string
	Filename     string
	Size         int
	Created_at   string
	Updated_at   string
	Download_url string
	Owner        SUser
}

type SHomework struct {
	Course_id         string
	Number            int
	Title             string
	Created_at        string
	Updated_at        string
	Begin_at          string
	Due_at            string
	Submissions_count int
	Marks_count       int
}

type attachment struct {
	Filename     string
	Size         string
	Download_url string
}

type Homework struct {
	SHomework
	Max_score   string
	Body        string
	Attachments []attachment
}

type SSubmission struct {
	Course_id       string
	Homework_number int
	Student_id      string
	Created_at      string
	Updated_at      string
	Marked_at       string
	Score           string
	Max_score       string
}

type Submission struct {
	SSubmission
	Body                string
	Attachments         []attachment
	Comment             string
	Comment_attachments []attachment
}
