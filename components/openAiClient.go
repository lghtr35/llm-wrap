package components

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/lghtr35/llm-wrap/models"
)

type OpenAiClient struct {
	config models.VendorConfig
	client *http.Client
}

func NewOpenAiClient(config models.VendorConfig) *OpenAiClient {
	return &OpenAiClient{config: config, client: &http.Client{}}
}

func (c *OpenAiClient) GenerateText(payload models.Prompt) (string, error) {
	log.Println("Requesting Generate Text from OpenAI")
	actualPayload := models.OpenAIPrompt{
		Prompt:              payload,
		MaxCompletionTokens: 1024,
		Stream:              false,
	}
	buf, err := json.Marshal(actualPayload)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest("POST", c.config.Url, bytes.NewBuffer(buf))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+c.config.ApiKey)

	response, err := c.client.Do(request)
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
		return "", errors.New("OpenAiClient: response is empty")
	}

	return resp.Choices[0].Message.Content, nil
}

func (c *OpenAiClient) GenerateTextAsStream(payload models.Prompt) (io.Reader, error) {
	log.Println("Requesting Generate Text As Stream from OpenAI")
	actualPayload := models.OpenAIPrompt{
		Prompt:              payload,
		MaxCompletionTokens: 1024,
		Stream:              true,
	}
	buf, err := json.Marshal(actualPayload)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", c.config.Url, bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+c.config.ApiKey)

	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}

	return response.Body, nil
}

func (c *OpenAiClient) CountTokens(payload models.OpenAIPrompt) int {
	return 0
}
