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
	// Identifiers.
	Id             string `json:"id,omitempty"`
	Semester       string `json:"semester,omitempty"`
	CourseNumber   string `json:"course_number,omitempty"`
	CourseSequence string `json:"course_sequence,omitempty"`

	// Metadata.
	Name        string `json:"name,omitempty"`
	Credit      int    `json:"credit,omitempty"`
	Hour        int    `json:"hour,omitempty"`
	Description string `json:"description,omitempty"`

	// Time & location.
	Weeks       string `json:"weeks,omitempty"`
	DayOfWeek   int    `json:"day_of_week,omitempty"`
	PeriodOfDay int    `json:"period_of_day,omitempty"`
	Location    string `json:"location,omitempty"`

	// Staff.
	// TODO: Pointer slice?
	Teachers   []*User `json:"teachers,omitempty"`
	Assistants []*User `json:"assistants,omitempty"`
}

type Announcement struct {
	// Identifiers.
	Id       string `json:"id,omitempty"`
	CourseId string `json:"course_id,omitempty"`

	// Metadata.
	Owner     User   `json:"owner,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	Priority  int    `json:"priority,omitempty"`

	// Content.
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
}

type File struct {
	// Identifiers.
	Id       string `json:"id,omitempty"`
	CourseId string `json:"course_id,omitempty"`

	// Metadata.
	Owner       User     `json:"owner,omitempty"`
	CreatedAt   string   `json:"created_at,omitempty"`
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Category    []string `json:"category,omitempty"`

	// Content.
	Filename    string `json:"filename,omitempty"`
	Size        int    `json:"size,omitempty"`
	DownloadUrl string `json:"download_url,omitempty"`
}

type Attachment struct {
	Filename    string `json:"filename,omitempty"`
	Size        int    `json:"size,omitempty"`
	DownloadUrl string `json:"download_url,omitempty"`
}

type Homework struct {
	// Identifiers.
	Id       string `json:"id,omitempty"`
	CourseId string `json:"course_id,omitempty"`

	// Metadata.
	Owner     User   `json:"owner,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	BeginAt   string `json:"begin_at,omitempty"`
	DueAt     string `json:"due_at,omitempty"`

	// Content.
	Title      string     `json:"title,omitempty"`
	Body       string     `json:"body,omitempty"`
	Attachment Attachment `json:"attachment,omitempty"`

	// Submissions.
	Submissions []*Submission `json:"submissions,omitempty"`
}

type Submission struct {
	// Metadata.
	Owner     User   `json:"student,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`

	// Content.
	Body       string     `json:"body,omitempty"`
	Attachment Attachment `json:"attachment,omitempty"`

	// Scoring metadata.
	ScoredBy User   `json:"scored_by,omitempty"`
	ScoredAt string `json:"scored_at,omitempty"`

	// Scoring content.
	Score             string     `json:"score,omitempty"`
	Comment           string     `json:"comment,omitempty"`
	CommentAttachment Attachment `json:"comment_attachment,omitempty"`
}
