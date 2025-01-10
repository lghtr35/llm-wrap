package components

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/lghtr35/llm-wrapping/models"
)

type LlmClient struct {
	config models.VendorConfig
}

func NewLlmClient(config models.VendorConfig) *LlmClient {
	return &LlmClient{config: config}
}

func (c *LlmClient) GenerateText(payload models.OpenAIPrompt) (string, error) {
	buf, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	response, err := http.Post(c.config.Url, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return "", err
	}

	str, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var resp models.OpenAIResponse
	err = json.Unmarshal(str, &resp)
	if err != nil {
		return "", err
	}
	if len(resp.Choices) < 1 {
		return "", errors.New("llmClient: response is empty")
	}

	return resp.Choices[0].Message.Content, nil
}

func (c *LlmClient) GenerateTextAsStream(payload models.OpenAIPrompt) (io.Reader, error) {
	buf, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	response, err := http.Post(c.config.Url, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}

	return response.Body, nil
}

func (c *LlmClient) CountTokens(payload models.OpenAIPrompt) int {
	return 0
}
