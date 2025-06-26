package service

import (
	// "strings"
	"testing"

	// "github.com/aziz-shoko/goblog/models"
	"github.com/aziz-shoko/goblog/internal/store"
)

func TestPostService_CreatePost(t *testing.T) {
	// t.Run("successfully creates post with valid data", func(t *testing.T) {
	// 	// setup
	// 	mockStore := store.NewInMemoryStore()
	// 	service := NewPostService(mockStore)

	// 	// Act
	// 	post, err := service.CreatePost("My first blog", "Some content for first blog")

	// 	// Assert
	// 	AssertError(t, err, nil)
	// 	AssertTest(t, post.Name, "My first blog")

	// 	// Verify it was actually stored
	// 	storedPost, err := mockStore.GetByID(post.ID)
	// 	AssertError(t, err, nil)
	// 	AssertTest(t, storedPost.Name, post.Name)
	// })

	// t.Run("sanitize title by trimming whitespace", func(t *testing.T) {
	// 	// Arrange
	// 	mockStore := store.NewInMemoryStore()
	// 	service := NewPostService(mockStore)

	// 	// Act - title with extra whitespace
	// 	post, err := service.CreatePost("  My Blog Post  ", "Some good content")
	// 	AssertError(t, err, nil)
	// 	AssertTest(t, post.Name, "My Blog Post")
	// })

	// t.Run("reject content that is too short", func(t *testing.T) {
	// 	// arrange
	// 	mockStore := store.NewInMemoryStore()
	// 	service := NewPostService(mockStore)

	// 	// Act - very short content
	// 	_, err := service.CreatePost("Good title", "Hi")
	// 	AssertError(t, err, ErrContentTooShort)
	// })

	// t.Run("prevent duplicate titles", func(t *testing.T) {
	// 	// arrange
	// 	mockStore := store.NewInMemoryStore()
	// 	service := NewPostService(mockStore)

	// 	// Act 1 - create first post
	// 	_, err := service.CreatePost("Unique Title", "First post content here")
	// 	AssertError(t, err, nil)

	// 	// Act 2 - try to create post with same title
	// 	_, err = service.CreatePost("Unique Title", "Second post content here")
	// 	AssertError(t, err, ErrDuplicateTitle)
	// })
	tests := []struct {
		name      string
		title     string
		content   string
		wantErr   error
		wantTitle string
	}{
		{
			name:      "successfully creates post with valid data",
			title:     "My first blog",
			content:   "Some content for first blog",
			wantErr:   nil,
			wantTitle: "My first blog",
		},
		{
			name:      "sanitize title by trimming whitespace",
			title:     "  My Blog Post  ",
			content:   "Some good content",
			wantErr:   nil,
			wantTitle: "My Blog Post",
		},
		{
			name:    "reject content that is too short",
			title:   "Good title",
			content: "Hi",
			wantErr: ErrContentTooShort,
		},
		{
			name:    "prevent duplicate titles",
			title:   "Unique Title",
			content: "Second post content here",
			wantErr: ErrDuplicateTitle,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup a fresh store for each test
			mockStore := store.NewInMemoryStore()
			service := NewPostService(mockStore)

			// For the dplicate-title case
			if tc.name == "prevent duplicate titles" {
				_, err := service.CreatePost(tc.title, "First post content here")
				if err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			}

			// Act
			post, err := service.CreatePost(tc.title, tc.content)

			// Assert error
			AssertError(t, err, tc.wantErr)

			if tc.wantErr == nil {
				AssertTest(t, post.Name, tc.wantTitle)
				stored, err := mockStore.GetByID(post.ID)
				AssertError(t, err, nil)
				AssertTest(t, stored.Name, tc.wantTitle)
			}
		})
	}

}

func AssertTest(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("Got %q, want %q", got, want)
	}
}

func AssertError(t testing.TB, got, want error) {
	t.Helper()
	if want == nil && got != nil {
		t.Fatalf("Expected no error but got %v", got)
	}

	if got != want {
		t.Errorf("Got error %v want error %v", got, want)
	}
}
