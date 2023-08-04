package models

import "time"

type ContentHistory struct {
	ID            string    `json:"_id"`
	ContentId     string    `json:"content_id"`
	UserId        string    `json:"user_id"`
	Action        string    `json:"action"`
	PreviousValue string    `json:"previous_value"`
	NewValue      string    `json:"new_value"`
	CreatedAt     time.Time `json:"created_at"`
}
