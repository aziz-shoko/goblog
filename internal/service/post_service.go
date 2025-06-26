package service

import (
	"errors"
	"strings"

	"github.com/aziz-shoko/goblog/internal/store"
	"github.com/aziz-shoko/goblog/models"
)

var (
	ErrContentTooShort = errors.New("Content Too Short, must be at least contain 5 chars")
	ErrDuplicateTitle = errors.New("Title already exists (case insensitive)")
)

// Post service handles business operations for blog posts
// Design pattern: Dependency Injection - depends on store interface
type PostServiceRepository struct {
	Store store.PostStore
}

// NewPostService creates a new post service
// Design pattern: Dependency Injection - inject the store dependency
func NewPostService(store store.PostStore) *PostServiceRepository {
	return &PostServiceRepository{
		Store: store,
	}
}

// CreatePost creates a new blog post with business rule validation
func (s *PostServiceRepository) CreatePost(title, content string) (*models.Post, error) {
	// Business rule 1: sanitize title
	trimmedTitle := strings.TrimSpace(title)

	// Business rule 2: validate the title and content
	if len(content) < 5 {
		return nil, ErrContentTooShort
	}

	// Business rule 3
	if s.titleExists(title) {
		return nil, ErrDuplicateTitle
	}

	// Create the post (using domain validation)
	post, err := models.NewPost(trimmedTitle, content)
	if err != nil {
		return nil, err
	}

	// store the post
	err = s.Store.Create(post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostServiceRepository) GetPostByID(id string) (*models.Post, error) {
	post, err := s.Store.GetByID(id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostServiceRepository) titleExists(title string) bool {
	posts, err := s.Store.GetAll()
	if err != nil {
		return false // if we cant check, return false
	}

	for _, post := range posts {
		if strings.EqualFold(post.Name, title) {
			return true
		}
	}

	return false
}
