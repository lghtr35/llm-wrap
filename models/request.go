package models

import "github.com/invopop/jsonschema"

type CommandRequest struct {
	Prompt string `json:"prompt"`
}

type OpenAIPrompt struct {
	MaxCompletionTokens int             `json:"max_completion_tokens"`
	Stream              bool            `json:"stream"`
	ResponseFormat      *ResponseFormat `json:"response_format"`
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

type ResponseFormat struct {
	Type       string     `json:"type"`
	JsonSchema JsonSchema `json:"json_schema"`
}

type JsonSchema struct {
	Name   string `json:"name"`
	Schema any    `json:"schema"`
}

func NewResponseFormat(schema *jsonschema.Schema) (ResponseFormat, error) {
	return ResponseFormat{
		Type: "json_schema",
		JsonSchema: JsonSchema{
			Name:   "response",
			Schema: schema,
		},
	}, nil
}
