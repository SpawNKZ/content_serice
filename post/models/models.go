package models

type Post struct {
	ID          string   `json:"id"`
	Category    string   `json:"category"`
	Resources   []string `json:"resources"`
	ContentId   string   `json:"content_id"`
	Description string   `json:"description"`
}

type Pagination struct {
	Total      int  `json:"total"`
	Limit      int  `json:"limit"`
	Offset     int  `json:"offset"`
	IsLastPage bool `json:"isLastPage"`
}

const (
	DefaultLimit = 10
	MaxOffset    = 2147483647
	MaxLimit     = 1000
)
