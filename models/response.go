package models

type OpenAIResponse struct {
	ID      string                 `json:"id"`
	Model   string                 `json:"model"`
	Choices []OpenAIResponseChoice `json:"choices"`
}

type OpenAIResponseChoice struct {
	Message      OpenAIResponseMessage `json:"message"`
	FinishReason *string               `json:"finish_reason"`
}

type OpenAIResponseMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
