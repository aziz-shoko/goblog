package models

import (
	"errors"
	"time"
	"github.com/google/uuid"
)

var (
	ErrEmtpyTitle   = errors.New("title cannot be emtpy")
	ErrEmtpyContent = errors.New("content cannot be empty")
)

type Post struct {
	Name      string
	Content   string
	ID        string
	CreatedAt time.Time
}

func NewPost(name, content string) (*Post, error) {
	if name == "" {
		return nil, ErrEmtpyTitle
	} else if content == "" {
		return nil, ErrEmtpyContent
	}

	return &Post{
		Name:      name,
		Content:   content,
		ID:        uuid.NewString(),
		CreatedAt: time.Now().UTC(),
	}, nil
}
