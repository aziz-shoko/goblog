package models

import (
	"testing"
	"time"
)

func TestNewPost_Validation(t *testing.T) {
	cases := []struct {
		Name    string
		Title   string
		Content string
		WantErr error
	}{
		{Name: "emtpy title test", Title: "", Content: "some content", WantErr: ErrEmtpyTitle},
		{Name: "emtpy content test", Title: "test", Content: "", WantErr: ErrEmtpyContent},
	}

	for _, tc := range cases {
		_, err := NewPost(tc.Title, tc.Content)
		AssertError(t, err, tc.WantErr)
	}
}

// Side affects - basically anything that a function does besides returning a value
// For example: Modifying variables, I/O operations such as printing, writing to files, or network reqs
func TestNewPost_SideAffects(t *testing.T) {
	t.Run("Set ID and CreatedAt", func(t *testing.T) {
		post, err := NewPost("title", "content")
		if err != nil {
			t.Fatal(err)
		}

		if post.ID == "" {
			t.Errorf("Expected non-empty ID")
		}

		if time.Since(post.CreatedAt) > time.Second {
			t.Errorf("CreatedAt should be recent")
		}
	})
}

func AssertError(t testing.TB, got, want error) {
	t.Helper()
	if got != want {
		t.Errorf("Got %q wanted %q", got, want)
	}
}
