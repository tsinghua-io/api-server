package model

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

type TimeLocation struct {
	Weeks       string `json:"weeks,omitempty"`
	DayOfWeek   int    `json:"day_of_week,omitempty"`
	PeriodOfDay int    `json:"period_of_day,omitempty"`
	Location    string `json:"location,omitempty"`
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
	TimeLocations []*TimeLocation `json:"time_locations,omitempty"`

	// Staff.
	Teachers   []*User `json:"teachers,omitempty"`
	Assistants []*User `json:"assistants,omitempty"`
}

type Announcement struct {
	// Identifiers.
	Id       string `json:"id,omitempty"`
	CourseId string `json:"course_id,omitempty"`

	// Metadata.
	Owner     *User  `json:"owner,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	Priority  int    `json:"priority"`

	// Content.
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
}

type File struct {
	// Identifiers.
	Id       string `json:"id,omitempty"`
	CourseId string `json:"course_id,omitempty"`

	// Metadata.
	Owner       *User    `json:"owner,omitempty"`
	CreatedAt   string   `json:"created_at,omitempty"`
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Category    []string `json:"category,omitempty"`

	// Content.
	Filename    string `json:"filename,omitempty"`
	Size        int    `json:"size"`
	DownloadURL string `json:"download_url,omitempty"`
}

type Attachment struct {
	Filename    string `json:"filename,omitempty"`
	Size        int    `json:"size"`
	DownloadURL string `json:"download_url,omitempty"`
}

type Submission struct {
	// Metadata.
	Owner        *User  `json:"student,omitempty"`
	AssignmentId string `json:"assignment_id,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	Late         bool   `json:"late"`

	// Content.
	Body       string      `json:"body,omitempty"`
	Attachment *Attachment `json:"attachment,omitempty"`

	// Scoring metadata.
	MarkedBy *User  `json:"marked_by,omitempty"`
	MarkedAt string `json:"marked_at,omitempty"`

	// Scoring content.
	Mark              *float32    `json:"mark"`
	Comment           string      `json:"comment,omitempty"`
	CommentAttachment *Attachment `json:"comment_attachment,omitempty"`
}

type Assignment struct {
	// Identifiers.
	Id       string `json:"id,omitempty"`
	CourseId string `json:"course_id,omitempty"`

	// Metadata.
	CreatedAt string `json:"created_at,omitempty"`
	BeginAt   string `json:"begin_at,omitempty"`
	DueAt     string `json:"due_at,omitempty"`

	// Content.
	Title      string      `json:"title,omitempty"`
	Body       string      `json:"body,omitempty"`
	Attachment *Attachment `json:"attachment,omitempty"`

	// Submissions.
	Submissions []*Submission `json:"submissions,omitempty"`
}
