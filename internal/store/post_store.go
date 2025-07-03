package store

import (
	"github.com/aziz-shoko/goblog/models"
	"errors"
)

var (
	ErrNotFound = errors.New("Item not found")
)

type InMemoryStore struct {
	posts map[string]*models.Post
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		posts: make(map[string]*models.Post),
	}
}

func (s *InMemoryStore) Create(post *models.Post) error {
	if post == nil {
		return errors.New("post cannot be nil")
	}

	s.posts[post.ID] = post

	return nil	
}

func (s *InMemoryStore) GetByID(id string) (*models.Post, error) {
	if _, ok := s.posts[id]; !ok {
		return nil, ErrNotFound
	}
	return s.posts[id], nil
}

func (s *InMemoryStore) GetAll() ([]*models.Post, error) {
	if len(s.posts) == 0 {
		return nil, errors.New("Emtpy store")
	}
	listOfPosts := []*models.Post{}
	for _, val := range s.posts {
		listOfPosts = append(listOfPosts, val)
	}
	return listOfPosts, nil
}

func (s *InMemoryStore) DeleteAll() error {
	s.posts = make(map[string]*models.Post)	
	return nil
}