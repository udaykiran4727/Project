package models

import "time"

type Link struct {
	ID          int64     `json:"id"`
	Shortcut    string    `json:"shortcut"`
	Destination string    `json:"destination"`
	CreatedAt   time.Time `json:"created_at"`
	ClickCount  int64     `json:"click_count"`
}

type CreateLinkRequest struct {
	Shortcut    string `json:"shortcut"`
	Destination string `json:"destination"`
}
