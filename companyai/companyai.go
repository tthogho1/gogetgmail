package companyai

import (
	"context"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

// GetCompanyName extracts a company name from the given text using OpenAI Chat API.
// It returns a single company name string with no extra punctuation or explanation.
func GetCompanyName(ctx context.Context, client *openai.Client, text string) (string, error) {
	prompt := fmt.Sprintf(
		"以下の文章から会社名（法人名）だけを抽出してください。\n"+
			"会社名だけを1つだけ返してください。余計な説明や句読点は不要です。\n\n%s", text)

	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	})
	if err != nil {
		return "", fmt.Errorf("openai chat completion failed: %w", err)
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("openai returned no choices")
	}
	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}
