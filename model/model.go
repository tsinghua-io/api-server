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

type Schedule struct {
	Weeks    string `json:"weeks,omitempty"`
	Day      int    `json:"day,omitempty"`
	Slot     int    `json:"slot,omitempty"`
	Location string `json:"location,omitempty"`
}

type Course struct {
	// Identifiers.
	Id         string `json:"id,omitempty"`
	SemesterId string `json:"semester_id,omitempty"`
	Number     string `json:"number,omitempty"`
	Sequence   string `json:"sequence,omitempty"`

	// Metadata.
	Name        string `json:"name,omitempty"`
	Credit      int    `json:"credit,omitempty"`
	Hour        int    `json:"hour,omitempty"`
	Description string `json:"description,omitempty"`

	// Schedules.
	Schedules []*Schedule `json:"schedules,omitempty"`

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
	Owner        *User  `json:"owner,omitempty"`
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

	// Submission.
	Submission *Submission `json:"submission,omitempty"`
}

type Materials struct {
	Announcements []*Announcement `json:"announcements,omitempty"`
	Files         []*File         `json:"files,omitempty"`
	Assignments   []*Assignment   `json:"assignments,omitempty"`
}

type Semester struct {
	Id      string `json:"id,omitempty"`
	BeginAt string `json:"begin_at,omitempty"`
}
