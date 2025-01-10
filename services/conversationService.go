package services

import (
	"encoding/json"
	"errors"
	"strings"
	"sync"

	"github.com/lghtr35/llm-wrapping/components"
	"github.com/lghtr35/llm-wrapping/models"
)

type ConversationService struct {
	promptGenerator *components.PromptGenerator
	// Can be an array of LlmClient too but for simplicity I am using 3 different clients
	llm1 *components.LlmClient
	llm2 *components.LlmClient
	llm3 *components.LlmClient
}

func NewConversationService(vendorConfig models.VendorConfig) *ConversationService {
	return &ConversationService{
		promptGenerator: components.NewPromptGenerator(vendorConfig),
		// can be initialized with different vendorConfigs but for simplicity I am using the same config
		llm1: components.NewLlmClient(vendorConfig),
		llm2: components.NewLlmClient(vendorConfig),
		llm3: components.NewLlmClient(vendorConfig),
	}
}

func (s *ConversationService) GenerateConversation(convoChannel ConversationChannel) error {
	fullAnalysis := ""
	summary := ""
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go s.getSentimentAnalysis(convoChannel, &fullAnalysis, wg)
	go s.getSummary(convoChannel, &summary, wg)
	wg.Wait()
	if fullAnalysis == "" || summary == "" {
		return errors.New("conversationService: fullAnalysis or summary is empty")
	}
	convoChannel.StatusUpdate(models.COMBINING_ANSWERS)
	temp := strings.Split(fullAnalysis, "\n")
	analysis := temp[0]
	purpose := temp[1]

	complexPrompt := s.promptGenerator.GenerateComplexPrompt(analysis, purpose, summary)
	openAiPrompt := models.OpenAIPrompt{
		Prompt:              complexPrompt,
		MaxCompletionTokens: 1024,
		Stream:              true,
	}
	convoChannel.StatusUpdate(models.LLM3_STARTED)
	reader, err := s.llm3.GenerateTextAsStream(openAiPrompt)
	if err != nil {
		convoChannel.FullUpdate(models.ERROR, err.Error())
		return err
	}
	decoder := json.NewDecoder(reader)
	result := ""
	for {
		var message models.OpenAIResponse
		err := decoder.Decode(&message)
		if err != nil || message.Choices == nil || len(message.Choices) < 1 || message.Choices[0].FinishReason != nil {
			break
		}
		result += message.Choices[0].Message.Content
		convoChannel.FullUpdate(models.LLM3_STREAMING, result)
	}

	convoChannel.FullUpdate(models.LLM3_FINISHED, result)
	return nil
}

func (s *ConversationService) getSentimentAnalysis(convoChannel ConversationChannel, result *string, wg *sync.WaitGroup) {
	defer wg.Done()
	payload := convoChannel.Payload()
	sentimentAnalysisPrompt := s.promptGenerator.GenerateSentimentAnalysisPrompt(payload)
	openAiPrompt := models.OpenAIPrompt{
		Prompt:              sentimentAnalysisPrompt,
		MaxCompletionTokens: 1024,
		Stream:              false,
	}
	convoChannel.StatusUpdate(models.LLM1_STARTED)
	var err error
	*result, err = s.llm1.GenerateText(openAiPrompt)
	if err != nil {
		convoChannel.FullUpdate(models.ERROR, err.Error())
		return
	}
	convoChannel.StatusUpdate(models.LLM1_FINISHED)
}

func (s *ConversationService) getSummary(convoChannel ConversationChannel, result *string, wg *sync.WaitGroup) {
	defer wg.Done()
	payload := convoChannel.Payload()
	summaryPrompt := s.promptGenerator.GenerateSummarizePrompt(payload)
	openAiPrompt := models.OpenAIPrompt{
		Prompt:              summaryPrompt,
		MaxCompletionTokens: 1024,
		Stream:              false,
	}
	convoChannel.StatusUpdate(models.LLM2_STARTED)
	var err error
	*result, err = s.llm2.GenerateText(openAiPrompt)
	if err != nil {
		convoChannel.FullUpdate(models.ERROR, err.Error())
		return
	}
	convoChannel.StatusUpdate(models.LLM2_FINISHED)
}
