package models

type CommandRequest struct {
	Prompt string `json:"prompt"`
}

type OpenAIPrompt struct {
	MaxCompletionTokens int  `json:"max_completion_tokents"`
	Stream              bool `json:"stream"`
	Prompt
}

type Prompt struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
