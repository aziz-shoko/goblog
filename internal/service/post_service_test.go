package service

import (
	// "strings"
	"testing"

	// "github.com/aziz-shoko/goblog/models"
	"github.com/aziz-shoko/goblog/internal/store"
)

func TestPostService_CreatePost(t *testing.T) {
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

func TestPostService_Get(t *testing.T) {
	// Setup
	mockStore := store.NewInMemoryStore()
	service := NewPostService(mockStore)

	post, err := service.CreatePost("Get Test Title", "Test content for get")
	AssertError(t, err, nil)
	
	// test
	_, err = service.GetPostByID(post.ID)
	if err == store.ErrNotFound {
		t.Errorf("GetByID operation failed")
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
