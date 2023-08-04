package models

type ContentStatus struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	IsRemovable bool   `json:"is_removable"`
}
