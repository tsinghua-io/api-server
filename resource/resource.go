//Definitions of resources.

package resource

type User struct {
	Id         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Type       string `json:"type,omitempty"`
	Department string `json:"department,omitempty"`
	Class      string `json:"class,omitempty"`
	Gender     string `json:"gender,omitempty"`
	Email      string `json:"email,omitempty"`
	Phone      string `json:"phone,omitempty"`
}

type Course struct {
	Id             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Teacher        User   `json:"teacher,omitempty"`
	Coteachers     []User `json:"coteachers,omitempty"`
	SchoolYear     string `json:"school_year,omitempty"`
	Semester       string `json:"semester,omitempty"`
	CourseNumber   string `json:"course_number,omitempty"`
	CourseSequence string `json:"course_sequence,omitempty"`
	Credit         int    `json:"credit,omitempty"`
	Hour           int    `json:"hour,omitempty"`
	Description    string `json:"description,omitempty"`
	StudentCount   int    `json:"student_count,omitempty"`
}

type Announcement struct {
	Id        string `json:"id,omitempty"`
	CourseId  string `json:"course_id,omitempty"`
	Title     string `json:"title,omitempty"`
	Owner     User   `json:"owner,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	Important bool   `json:"important,omitempty"`
	Body      string `json:"body,omitempty"`
}

type File struct {
	Id          string   `json:"id,omitempty"`
	CourseId    string   `json:"course_id,omitempty"`
	Category    []string `json:"category,omitempty"`
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Filename    string   `json:"filename,omitempty"`
	Size        int      `json:"size,omitempty"`
	DownloadUrl string   `json:"download_url,omitempty"`
	CreatedAt   string   `json:"created_at,omitempty"`
	Owner       User     `json:"owner,omitempty"`
}

type Attachment struct {
	Filename    string `json:"filename,omitempty"`
	Size        int    `json:"size,omitempty"`
	DownloadUrl string `json:"download_url,omitempty"`
}

type Homework struct {
	Id              string     `json:"id,omitempty"`
	CourseId        string     `json:"course_id,omitempty"`
	Title           string     `json:"title,omitempty"`
	CreatedAt       string     `json:"created_at,omitempty"`
	BeginAt         string     `json:"begin_at,omitempty"`
	DueAt           string     `json:"due_at,omitempty"`
	SubmissionCount int        `json:"submission_count,omitempty"`
	MarkCount       int        `json:"mark_count,omitempty"`
	Body            string     `json:"body,omitempty"`
	Attachment      Attachment `json:"attachment,omitempty"`
}

type Submission struct {
	CourseId          string     `json:"course_id,omitempty"`
	HomeworkId        string     `json:"homework_id,omitempty"`
	Student           User       `json:"student,omitempty"`
	CreatedAt         string     `json:"created_at,omitempty"`
	MarkedAt          string     `json:"marked_at,omitempty"`
	Score             string     `json:"score,omitempty"`
	Body              string     `json:"body,omitempty"`
	Attachment        Attachment `json:"attachment,omitempty"`
	Comment           string     `json:"comment,omitempty"`
	CommentAttachment Attachment `json:"comment_attachment,omitempty"`
}
