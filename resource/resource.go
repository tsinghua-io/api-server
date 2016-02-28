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

type TimeLocation struct {
	Weeks       string `json:"weeks,omitempty"`
	DayOfWeek   int    `json:"day_of_week"`
	PeriodOfDay int    `json:"period_of_day"`
	Location    string `json:"location,omitempty"`
}

type Course struct {
	// Identifiers.
	Id             string `json:"id"`
	Semester       string `json:"semester"`
	CourseNumber   string `json:"course_number"`
	CourseSequence string `json:"course_sequence"`

	// Metadata.
	Name        string `json:"name"`
	Credit      int    `json:"credit"`
	Hour        int    `json:"hour"`
	Description string `json:"description,omitempty"`

	// Time & location.
	TimeLocations []*TimeLocation `json:"time_locations,omitempty"`

	// Staff.
	Teachers   []*User `json:"teachers,omitempty"`
	Assistants []*User `json:"assistants,omitempty"`
}

type Announcement struct {
	// Identifiers.
	Id       string `json:"id"`
	CourseId string `json:"course_id"`

	// Metadata.
	Owner     *User  `json:"owner,omitempty"`
	CreatedAt string `json:"created_at"`
	Priority  int    `json:"priority"`
	Read      bool   `json:"read"`

	// Content.
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
}

type File struct {
	// Identifiers.
	Id       string `json:"id"`
	CourseId string `json:"course_id"`

	// Metadata.
	Owner       *User    `json:"owner,omitempty"`
	CreatedAt   string   `json:"created_at"`
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Category    []string `json:"category,omitempty"`
	Read        bool     `json:"read"`

	// Content.
	Filename    string `json:"filename"`
	Size        int    `json:"size"`
	DownloadUrl string `json:"download_url"`
}

type Attachment struct {
	Filename    string `json:"filename"`
	Size        int    `json:"size"`
	DownloadUrl string `json:"download_url"`
}

type Submission struct {
	// Metadata.
	Owner     *User  `json:"student,omitempty"`
	CreatedAt string `json:"created_at"`
	Late      bool   `json:"late"`

	// Content.
	Body       string      `json:"body,omitempty"`
	Attachment *Attachment `json:"attachment,omitempty"`

	// Scoring metadata.
	MarkedBy *User  `json:"marked_by,omitempty"`
	MarkedAt string `json:"marked_at,omitempty"`

	// Scoring content.
	Mark              *float32    `json:"mark,omitempty"`
	Comment           string      `json:"comment,omitempty"`
	CommentAttachment *Attachment `json:"comment_attachment,omitempty"`
}

type Homework struct {
	// Identifiers.
	Id       string `json:"id"`
	CourseId string `json:"course_id"`

	// Metadata.
	CreatedAt         string `json:"created_at,omitempty"`
	BeginAt           string `json:"begin_at"`
	DueAt             string `json:"due_at"`
	SubmittedCount    int    `json:"submitted_count,omitempty"`
	NotSubmittedCount int    `json:"not_submitted_count,omitempty"`
	SeenCount         int    `json:"seen_count,omitempty"`
	MarkedCount       int    `json:"marked_count,omitempty"`

	// Content.
	Title      string      `json:"title,omitempty"`
	Body       string      `json:"body,omitempty"`
	Attachment *Attachment `json:"attachment,omitempty"`

	// Submissions.
	Submissions []*Submission `json:"submissions,omitempty"`
}
