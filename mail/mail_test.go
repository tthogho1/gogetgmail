package mail

import (
	"encoding/base64"
	"testing"

	"google.golang.org/api/gmail/v1"
)

func TestGetHeader(t *testing.T) {
	msg := &gmail.Message{
		Payload: &gmail.MessagePart{
			Headers: []*gmail.MessagePartHeader{
				{Name: "Subject", Value: "Hello"},
				{Name: "From", Value: "me@example.com"},
			},
		},
	}
	if got := GetHeader(msg, "Subject"); got != "Hello" {
		t.Fatalf("expected Subject Hello, got %q", got)
	}
	if got := GetHeader(msg, "From"); got != "me@example.com" {
		t.Fatalf("expected From me@example.com, got %q", got)
	}
}

func TestGetBodyText_Simple(t *testing.T) {
	body := "simple body"
	enc := base64.URLEncoding.EncodeToString([]byte(body))
	msg := &gmail.Message{
		Payload: &gmail.MessagePart{
			MimeType: "text/plain",
			Body:     &gmail.MessagePartBody{Data: enc},
		},
	}
	got := GetBodyText(msg.Payload)
	if got != body {
		t.Fatalf("expected body %q, got %q", body, got)
	}
}

func TestGetBodyText_Recursive(t *testing.T) {
	body := "nested body"
	enc := base64.URLEncoding.EncodeToString([]byte(body))
	root := &gmail.MessagePart{
		MimeType: "multipart/alternative",
		Parts: []*gmail.MessagePart{
			{MimeType: "text/html", Body: &gmail.MessagePartBody{Data: ""}},
			{MimeType: "text/plain", Body: &gmail.MessagePartBody{Data: enc}},
		},
	}
	got := GetBodyText(root)
	if got != body {
		t.Fatalf("expected nested body %q, got %q", body, got)
	}
}
