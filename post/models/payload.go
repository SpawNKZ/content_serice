package models

type CreateRequest struct {
	Category    string   `json:"category"`
	Resources   []string `json:"resources"`
	ContentId   string   `json:"content_id"`
	Description string   `json:"description"`
}

type IdRequest struct {
	ID string
}

type UpdateRequest struct {
	ID          string
	Description string   `json:"description"`
	Resources   []string `json:"resources"`
}

type ReqPagination struct {
	Limit  int
	Offset int
}

type PostFilter struct {
	ContentId string `bson:"content_id"`
}

type GetListRequest struct {
	ReqPagination
	PostFilter
}
