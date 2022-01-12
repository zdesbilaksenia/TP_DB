package models

import "time"

type Post struct {
	Id       int       `json:"id,omitempty"`
	Parent   int       `json:"parent"`
	Author   string    `json:"author"`
	Message  string    `json:"message"`
	IsEdited bool      `json:"isEdited"`
	Forum    string    `json:"forum"`
	Thread   int       `json:"thread"`
	Created  time.Time `json:"created,omitempty"`
}

type Posts []*Post
