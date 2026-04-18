package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"

	"gogetgmail/companyai"
	"gogetgmail/mail"
	"gogetgmail/token"
)

func main() {
	// Load .env for OPENAI_API_KEY
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}
	aiClient := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	labelName := "INBOX"
	if len(os.Args) > 1 {
		labelName = os.Args[1]
	}

	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read credentials.json: %v", err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := token.GetClient(config)

	srv, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create Gmail service: %v", err)
	}

	labelID, err := mail.GetLabelID(srv, "me", labelName)
	if err != nil {
		log.Fatalf("Failed to find label: %v", err)
	}
	fmt.Printf("Label %q -> ID: %s\n\n", labelName, labelID)

	var allMessages []*gmail.Message
	pageToken := ""
	for {
		req := srv.Users.Messages.List("me").LabelIds(labelID).MaxResults(100)
		if pageToken != "" {
			req = req.PageToken(pageToken)
		}
		r, err := req.Do()
		if err != nil {
			log.Fatalf("Unable to list messages: %v", err)
		}
		allMessages = append(allMessages, r.Messages...)
		if r.NextPageToken == "" {
			break
		}
		pageToken = r.NextPageToken
	}

	fmt.Printf("Found %d message(s)\n\n", len(allMessages))

	limit := 10
	if len(allMessages) < limit {
		limit = len(allMessages)
	}

	var companies []string

	for i, m := range allMessages[:limit] {
		msg, err := srv.Users.Messages.Get("me", m.Id).Format("full").Do()
		if err != nil {
			log.Printf("Unable to get message %s: %v", m.Id, err)
			continue
		}

		subject := mail.GetHeader(msg, "Subject")
		from := mail.GetHeader(msg, "From")
		body := mail.GetBodyText(msg.Payload)

		fmt.Printf("--- Message %d ---\n", i+1)
		fmt.Printf("From:    %s\n", from)
		fmt.Printf("Subject: %s\n", subject)
		if body != "" {
			preview := strings.ReplaceAll(body, "\r\n", "\n")
			if len(preview) > 200 {
				preview = preview[:200] + "..."
			}
			fmt.Printf("Body:    %s\n", preview)

			// Extract company name via OpenAI
			company, err := companyai.GetCompanyName(context.Background(), aiClient, body)
			if err != nil {
				log.Printf("Company extraction failed: %v", err)
			} else {
				companies = append(companies, company)
			}
		}
		fmt.Println()
	}

	// Print all extracted company names
	fmt.Println("=== Extracted Companies ===")
	for i, c := range companies {
		fmt.Printf("  %d. %s\n", i+1, c)
	}
}
