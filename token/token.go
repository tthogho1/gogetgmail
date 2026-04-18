package token

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

// GetClient retrieves a token from file or web, saves it, and returns an HTTP client.
// It uses TokenSource to auto-refresh expired tokens.
func GetClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := FromFile(tokFile)
	if err != nil {
		tok = FromWeb(config)
		Save(tokFile, tok)
	}
	ts := config.TokenSource(context.Background(), tok)
	newTok, err := ts.Token()
	if err == nil && newTok.AccessToken != tok.AccessToken {
		Save(tokFile, newTok)
	}
	return oauth2.NewClient(context.Background(), ts)
}

// FromWeb requests a token from the web via user authorization.
func FromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code:\n%v\n", authURL)
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}
	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// FromFile retrieves a token from a local file.
func FromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Save saves a token to a file.
func Save(path string, t *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(t)
}
