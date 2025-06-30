package store

import (
	"fmt"
	"github.com/aziz-shoko/goblog/models"
	"strconv"
	"testing"
)

func TestPostStore_GetByID_Create(t *testing.T) {
	type scenario struct {
		name     string
		setup    func() *InMemoryStore
		post     *models.Post
		wantErr  bool
		validate func(s *InMemoryStore, post *models.Post) error
	}

	scenarios := []scenario{
		{
			name: "nil post returns error",
			setup: func() *InMemoryStore {
				return NewInMemoryStore()
			},
			post:    nil,
			wantErr: true,
		},
		{
			name: "successful create and retrieve",
			setup: func() *InMemoryStore {
				return NewInMemoryStore()
			},
			post: func() *models.Post {
				p, _ := models.NewPost("Test Title", "Test Content")
				return p
			}(),
			wantErr: false,
			validate: func(s *InMemoryStore, post *models.Post) error {
				got, err := s.GetByID(post.ID)
				if err != nil {
					return fmt.Errorf("GetByID failed: %v", err)
				}

				if got.Name != post.Name || got.Content != post.Content {
					return fmt.Errorf("stored post mismatch: got %+v, want %+v", got, post)
				}
				return nil
			},
		},
		{
			name: "failing retrieve",
			setup: func() *InMemoryStore {
				return NewInMemoryStore()
			},
			post: func() *models.Post {
				p, _ := models.NewPost("Test Title", "Test Content")
				return p
			}(),
			// false because we are not expecting error until validate
			wantErr: false,
			validate: func(s *InMemoryStore, post *models.Post) error {
				_, err := s.GetByID("somerandomnonsenseid")
				if err == nil {
					return fmt.Errorf("GetByID failed")
				}
				if err != ErrNotFound {
					return fmt.Errorf("Got error %v wanted error %v", err, ErrNotFound)
				}
				return nil
			},
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			store := sc.setup()
			err := store.Create(sc.post)
			if sc.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !sc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if sc.validate != nil {
				if err := sc.validate(store, sc.post); err != nil {
					t.Error(err)
				}
			}
		})
	}
}

// TODO: in the future, make this also a table driven to test all related edge cases to GetAll
// i will do this when i implement an actual database instead of just an in memory one
func TestPostStore_Get_Delete(t *testing.T) {
	// database setup
	database := NewInMemoryStore()

	// make posts and create them
	for i := range 5 {
		title := "Title" + strconv.Itoa(i)
		content := "Test Content" + strconv.Itoa(i)
		post, _ := models.NewPost(title, content)
		database.Create(post)
	}

	t.Run("Valid Get all test", func(t *testing.T) {
		listOfPosts, _ := database.GetAll()
		if len(listOfPosts) != 5 {
			t.Errorf("Expected 5 posts, got %d", len(listOfPosts))
		}
	})

	t.Run("Delete all posts", func(t *testing.T) {
		err := database.DeleteAll()
		if err != nil {
			t.Fatalf("Expected to delete all posts but failed: %v", err)
		}

		if _, err := database.GetAll(); err == nil {
			t.Error("Failed to delete all posts in database")
		}

	})
}

