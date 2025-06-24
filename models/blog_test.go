package models

import (
	"testing"
)

func TestNewPostValidation(t *testing.T) {
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



func AssertError(t testing.TB, got, want error) {
	t.Helper()
	if got != want {
		t.Errorf("Got %q wanted %q", got, want)
	}
}
