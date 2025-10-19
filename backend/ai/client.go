package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

type AIClient interface {
	AnalyzeText(ctx context.Context, prompt string, models []string) (string, error)
}

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient() AIClient {
	return &Client{
		apiKey:     os.Getenv("GROQ_API_KEY"),
		httpClient: &http.Client{},
	}
}

type GroqRequest struct {
	Messages []Message `json:"messages"`
	Model    string    `json:"model"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GroqResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

const systemPrompt = "You are a supportive and knowledgeable health assistant. Your goal is to provide encouraging and safe wellness advice based on the user's logs. Be positive, focus on actionable tips, and keep your responses concise and to the point. Do not use emojis, Don't be corny nor cringe, Don't use markdown, keep it 2 short sentences MAX."

func (c *Client) AnalyzeText(ctx context.Context, prompt string, models []string) (string, error) {
	var lastErr error

	for _, model := range models {
		reqPayload := GroqRequest{
			Messages: []Message{
				{Role: "system", Content: systemPrompt},
				{Role: "user", Content: prompt},
			},
			Model: model,
		}

		reqBody, err := json.Marshal(reqPayload)
		if err != nil {
			return "", err
		}

		req, err := http.NewRequestWithContext(ctx, "POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(reqBody))
		if err != nil {
			return "", err
		}

		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		if resp.StatusCode == 429 || resp.StatusCode >= 500 {
			lastErr = errors.New("Groq API returned retryable status: " + resp.Status)
			resp.Body.Close()
			continue
		}

		if resp.StatusCode >= 400 {
			defer resp.Body.Close()
			return "", errors.New("Groq API returned non-retryable status: " + resp.Status)
		}

		var groqResp GroqResponse
		if err := json.NewDecoder(resp.Body).Decode(&groqResp); err != nil {
			lastErr = err
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		if len(groqResp.Choices) > 0 && groqResp.Choices[0].Message.Content != "" {
			return groqResp.Choices[0].Message.Content, nil
		}
	}

	if lastErr != nil {
		return "", lastErr
	}

	return "", errors.New("all models failed to generate a response")
}
