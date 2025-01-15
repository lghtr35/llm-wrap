package components

import (
	"fmt"

	"github.com/lghtr35/llm-wrap/models"
)

type PromptGenerator struct {
	vendorConfig models.VendorConfig
}

func NewPromptGenerator(vendorConfig models.VendorConfig) *PromptGenerator {
	return &PromptGenerator{vendorConfig: vendorConfig}
}

func (g *PromptGenerator) GenerateSentimentAnalysisPrompt(payload string) models.Prompt {
	return models.Prompt{
		Model: g.vendorConfig.Model,
		Messages: []models.Message{
			{
				Role:    g.getDeveloperRole(),
				Content: "Please give sentiment and purpose of the given message in a short and concise way without any further explanation",
			},
			{
				Role:    "user",
				Content: payload,
			},
		},
	}
}

func (g *PromptGenerator) GenerateSummarizePrompt(payload string) models.Prompt {
	return models.Prompt{
		Model: g.vendorConfig.Model,
		Messages: []models.Message{
			{
				Role:    g.getDeveloperRole(),
				Content: "Please shorten and summarize the given message withouth losing any core information and adding any extra explanation",
			},
			{
				Role:    "user",
				Content: payload,
			},
		},
	}
}

func (g *PromptGenerator) GenerateComplexPrompt(sentiment, purpose, summary string) models.Prompt {
	return models.Prompt{
		Model: g.vendorConfig.Model,
		Messages: []models.Message{
			{
				Role:    "user",
				Content: fmt.Sprintf("Please answer to this message taking into account that the sentiment is %s and the purpose is %s. Message: %s", sentiment, purpose, summary),
			},
		},
	}
}

func (g *PromptGenerator) getDeveloperRole() string {
	role := "developer"
	if g.vendorConfig.Name != "OpenAI" {
		role = "user"
	}
	return role
}
