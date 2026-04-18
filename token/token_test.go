package token

import (
	"os"
	"testing"

	"golang.org/x/oauth2"
)

func TestSaveAndFromFile(t *testing.T) {
	f, err := os.CreateTemp("", "tok-*.json")
	if err != nil {
		t.Fatalf("temp file: %v", err)
	}
	path := f.Name()
	f.Close()
	defer os.Remove(path)

	tok := &oauth2.Token{AccessToken: "AT", RefreshToken: "RT"}
	Save(path, tok)

	r, err := FromFile(path)
	if err != nil {
		t.Fatalf("FromFile error: %v", err)
	}
	if r.AccessToken != tok.AccessToken || r.RefreshToken != tok.RefreshToken {
		t.Fatalf("tokens differ: got %+v want %+v", r, tok)
	}
}
