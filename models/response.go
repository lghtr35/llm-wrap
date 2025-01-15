package models

import "github.com/invopop/jsonschema"

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

type OpenAIStreamingResponse struct {
	ID      string                  `json:"id"`
	Model   string                  `json:"model"`
	Choices []OpenAIStreamingChoice `json:"choices"`
}

type OpenAIStreamingChoice struct {
	Delta OpenAIResponseMessage `json:"delta"`
}

type OpenAIStreamingChoiceDelta struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type SentimentAnalysisResponse struct {
	Sentiment string `json:"sentiment"`
	Purpose   string `json:"purpose"`
}

func (s *SentimentAnalysisResponse) GetSchema() *jsonschema.Schema {
	schema := jsonschema.Reflect(&s)
	return schema.Definitions["SentimentAnalysisResponse"]
}

type SummaryAnalysisResponse struct {
	Summary string `json:"summary"`
}

func (s *SummaryAnalysisResponse) GetSchema() *jsonschema.Schema {
	schema := jsonschema.Reflect(&s)
	return schema.Definitions["SummaryAnalysisResponse"]
}

type OpenAIErrorResponse struct {
	Error InnerError `json:"error"`
}

type InnerError struct {
	Code    any    `json:"code,omitempty"`
	Message string `json:"message"`
}
