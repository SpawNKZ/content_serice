package models

import "time"

type Content struct {
	ID           string    `json:"id"`
	Locale       string    `json:"locale"`
	Body         string    `json:"body"`
	Description  string    `json:"description"`
	Resources    []string  `json:"resources"`
	SubjectId    int64     `json:"subject_id"`
	Subject      *Subject  `json:"subject"`
	MicrotopicId int64     `json:"microtopic_id"`
	StatusId     string    `json:"status_id"`
	AuthorId     string    `json:"author_id"`
	Difficulty   int       `json:"difficulty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type SubjectBase struct {
	ID        int64     `json:"id"`
	ProgramID int64     `json:"program_id"`
	ImageUrl  string    `json:"image_url"`
	IconUrl   string    `json:"icon_url"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Translation struct {
	ID     int64  `json:"translationId,omitempty"`
	Locale string `json:"locale" fake:"??"`
	Name   string `json:"name"`
}

type Subject struct {
	SubjectBase
	Translations []*Translation `json:"translations"`
}

type MicrotopicBase struct {
	ID          int64     `json:"id"`
	SubjectId   int64     `json:"subject_id"`
	SectionId   int64     `json:"section_id"`
	ObjectiveId int64     `json:"objective_id"`
	Quarter     int8      `json:"quarter"`
	GradeID     int64     `json:"grade_id"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Microtopic struct {
	MicrotopicBase
	Translations []*Translation `json:"translations"`
}
