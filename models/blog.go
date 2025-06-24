package models

import (
	"errors"
)

var (
	ErrEmtpyTitle = errors.New("title cannot be emtpy")
	ErrEmtpyContent = errors.New("content cannot be empty")
)

type Post struct {
	Name    string
	Content string
}

func NewPost(name, content string) (*Post, error) {
	if name == "" {
		return nil, ErrEmtpyTitle
	} else if content == "" {
		return nil, ErrEmtpyContent
	}

	return &Post{Name: name, Content: content}, nil
}
