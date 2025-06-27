package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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

