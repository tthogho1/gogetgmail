package mail

import (
	"encoding/base64"
	"fmt"

	"google.golang.org/api/gmail/v1"
)

// GetLabelID returns the label ID for a given label name.
func GetLabelID(srv *gmail.Service, user, labelName string) (string, error) {
	r, err := srv.Users.Labels.List(user).Do()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve labels: %v", err)
	}
	for _, l := range r.Labels {
		if l.Name == labelName {
			return l.Id, nil
		}
	}
	return "", fmt.Errorf("label %q not found", labelName)
}

// GetHeader extracts a header value by name from a Gmail message.
func GetHeader(msg *gmail.Message, name string) string {
	for _, h := range msg.Payload.Headers {
		if h.Name == name {
			return h.Value
		}
	}
	return ""
}

// GetBodyText extracts plain text body from a message, searching parts recursively.
func GetBodyText(payload *gmail.MessagePart) string {
	if payload.MimeType == "text/plain" && payload.Body != nil && payload.Body.Data != "" {
		data, err := base64.URLEncoding.DecodeString(payload.Body.Data)
		if err != nil {
			return ""
		}
		return string(data)
	}
	for _, part := range payload.Parts {
		if text := GetBodyText(part); text != "" {
			return text
		}
	}
	return ""
}
