package models

type CreateRequest struct {
	Locale       string   `json:"locale"`
	Body         string   `json:"body"`
	Description  string   `json:"description"`
	Resources    []string `json:"resources"`
	SubjectId    int64    `json:"subject_id"`
	MicrotopicId int64    `json:"microtopic_id"`
	StatusId     string
	AuthorId     string `json:"author_id"` //TODO: get from service
	Difficulty   int    `json:"difficulty"`
}

type IdRequest struct {
	ID string
}

type UpdateRequest struct {
	ID          string
	Body        string `json:"body"`
	Description string `json:"description"`
	Difficulty  int    `json:"difficulty"`
}

type AssignAuthorRequest struct {
	ID       string
	AuthorId string `json:"author_id"`
}

type ChangeStatusRequest struct {
	ID       string
	StatusId string `json:"status_id"`
}

type ReqPagination struct {
	Limit  int
	Offset int
}

type ContentFilter struct {
	Locale       string `bson:"locale"`
	Status       string `bson:"status"`
	SubjectId    int64  `bson:"subject_id"`
	MicrotopicId int64  `bson:"microtopic_id"`
	AuthorId     string `bson:"author_id"`
}

type GetListRequest struct {
	ReqPagination
	ContentFilter
}
