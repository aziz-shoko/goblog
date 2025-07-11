package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/aziz-shoko/goblog/internal/service"
	"github.com/aziz-shoko/goblog/internal/store"
	"github.com/aziz-shoko/goblog/models"
)

func TestPostHandler_CreatePost(t *testing.T) {

	store := store.NewInMemoryStore()
	service := service.NewPostService(store)
	handler := NewPostHandler(service)

	tests := []struct {
		name       string
		body       interface{}
		wantStatus int
		wantTitle  string
		wantError  bool
	}{
		{
			name:       "success",
			body:       CreatePostRequest{"My Name", "My content"},
			wantStatus: http.StatusCreated,
			wantTitle:  "My Name",
		},
		{
			name:       "bad json data",
			body:       []byte(`{"bad":}`), // raw invalid bytes
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
		{
			name:       "short content",
			body:       CreatePostRequest{"Test Name", "hi"},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w, req := setupTest(t, tc.body)
			handler.CreatePost(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("expected %d, got %d", tc.wantStatus, w.Code)
			}

			if tc.wantStatus == http.StatusCreated {
				var resp CreatePostResponse
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("unmarshal: %v", err)
				}
				if resp.Name != tc.wantTitle {
					t.Errorf("expected title %q, got %q", tc.wantTitle, resp.Name)
				}
			}
		})
	}
}

func setupTest(t *testing.T, body interface{}) (*httptest.ResponseRecorder, *http.Request) {
	t.Helper()

	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	return w, req
}

func TestLoggingMiddleware(t *testing.T) {
	t.Run("log request details", func(t *testing.T) {
		// Capture log output
		var logOutput bytes.Buffer
		log.SetOutput(&logOutput)
		defer log.SetOutput(os.Stdout) // Reset after test

		// Create a test handler
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("test response"))
		})

		// Wrap with logging middleware
		wrappedHandler := LoggingMiddleware(testHandler)

		// Make request
		req := httptest.NewRequest(http.MethodPost, "/posts", nil)
		w := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(w, req)

		// Assert log output
		logStr := logOutput.String()
		if !strings.Contains(logStr, "POST") {
			t.Error("Expected log to contain HTTP method")
		}
		if !strings.Contains(logStr, "/posts") {
			t.Error("Expected log to contain URL path")
		}
		if !strings.Contains(logStr, "201") {
			t.Error("Expected log to contain status code")
		}
	})
}

func TestPostHandler_Get(t *testing.T) {
	// setup
	store := store.NewInMemoryStore()
	service := service.NewPostService(store)
	handler := NewPostHandler(service)

	var response CreatePostResponse

	// make specific post for this test
	post, err := service.CreatePost("some test title", "some test content for this")
	if err != nil {
		t.Fatalf("Error creating posts")
	}

	tests := []struct {
		name       string
		wantStatus int
		wantID     string
		wantError  bool
	}{
		{
			name:       "Test Valid GetPostByID",
			wantStatus: http.StatusOK,
			wantID:     post.ID,
			wantError:  false,
		},
		{
			name:       "Test Invalid GetPostByID",
			wantStatus: http.StatusNotFound,
			wantID:     "invalidID", // Invalid id to test 404 not found
			wantError:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// make request
			req := httptest.NewRequest(http.MethodGet, "/post/"+tc.wantID, nil)
			w := httptest.NewRecorder()
			handler.GetPostByID(w, req)

			if w.Code != tc.wantStatus {
				t.Fatalf("expected status code %d but got %d", tc.wantStatus, w.Code)
			}

			if !tc.wantError {
				err = json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}

				if !tc.wantError && response.ID != post.ID {
					t.Errorf("expected %v got ID %v", post, response)
				}
			}
		})
	}
}

func TestPostHandler_GetAll(t *testing.T) {
	// setup
	store := store.NewInMemoryStore()
	service := service.NewPostService(store)
	handler := NewPostHandler(service)

	// Create some posts
	posts := []*models.Post{}
	for i := range 5 {
		post, err := service.CreatePost(strconv.Itoa(i)+"Name", "Some Content"+strconv.Itoa(i))
		if err != nil {
			t.Fatalf("Error creating posts")
		}
		posts = append(posts, post)
	}

	t.Run("Get All Posts", func(t *testing.T) {
		// Request
		req := httptest.NewRequest(http.MethodGet, "/posts/", nil)
		w := httptest.NewRecorder()
		handler.GetPostsAll(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, but got %d", http.StatusOK, w.Code)
		}

		// Check content type
		expectedContentType := "application/json"
		if contentType := w.Header().Get("Content-Type"); contentType != expectedContentType {
			t.Errorf("Expected Content-Type %s, got %s", expectedContentType, contentType)
		}

		var responseData []*models.Post
		err := json.Unmarshal(w.Body.Bytes(), &responseData)
		if err != nil {
			t.Fatalf("Error unmarshaling response: %v", err)
		}

		// Check that we got the expected number of posts
		// for simplicity sake, we will just check length instead of actual content
		// may god forgive me for this
		if len(responseData) != len(posts) {
			t.Fatalf("Expected %d posts, got %d", len(posts), len(responseData))
		}
	})

	t.Run("Delete All Posts", func(t *testing.T) {
		// Request 
		req := httptest.NewRequest(http.MethodDelete, "/posts", nil)
		w := httptest.NewRecorder()
		handler.DeleteAllPosts(w, req)

		if w.Code != http.StatusNoContent {
			t.Fatalf("Expected status code %d, but got %d", http.StatusNoContent, w.Code)
		}

		// test to see if content was actually deleted by called store method directly
		_, err := store.GetAll()
		if err == nil {
			t.Errorf("expected error for emtpy content but got nil")
		}

	})
}
