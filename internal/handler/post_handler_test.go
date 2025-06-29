package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"os"
	"strings"

	"github.com/aziz-shoko/goblog/internal/service"
	"github.com/aziz-shoko/goblog/internal/store"
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
			body:       CreatePostRequest{"My Title", "My content"},
			wantStatus: http.StatusCreated,
			wantTitle:  "My Title",
		},
		{
			name:       "bad json data",
			body:       []byte(`{"bad":}`), // raw invalid bytes
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
		{
			name:       "short content",
			body:       CreatePostRequest{"Test Title", "hi"},
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
				if resp.Title != tc.wantTitle {
					t.Errorf("expected title %q, got %q", tc.wantTitle, resp.Title)
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

