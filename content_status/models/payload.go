package models

type CreateRequest struct {
	Name        string `json:"name"`
	IsRemovable bool   `json:"is_removable"`
}

type IdRequest struct {
	ID string
}

type UpdateRequest struct {
	ID   string
	Name string `json:"name"`
}

type ReqPagination struct {
	Limit  int
	Offset int
}

type GetListRequest struct {
	ReqPagination
}
