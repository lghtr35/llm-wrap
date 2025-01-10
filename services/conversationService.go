package services

import (
	"encoding/json"
	"log"
	"strings"
	"sync"

	"github.com/lghtr35/llm-wrap/components"
	"github.com/lghtr35/llm-wrap/models"
)

type ConversationService struct {
	promptGenerator *components.PromptGenerator
	// Can be an array of OpenAiClient too but for simplicity I am using 3 different clients
	llm1 *components.OpenAiClient
	llm2 *components.OpenAiClient
	llm3 *components.OpenAiClient
}

func NewConversationService(vendorConfig models.VendorConfig) *ConversationService {
	return &ConversationService{
		promptGenerator: components.NewPromptGenerator(vendorConfig),
		// can be initialized with different vendorConfigs but for simplicity I am using the same config
		llm1: components.NewOpenAiClient(vendorConfig),
		llm2: components.NewOpenAiClient(vendorConfig),
		llm3: components.NewOpenAiClient(vendorConfig),
	}
}

func (s *ConversationService) GenerateConversation(convoChannel ConversationChannel) {
	defer convoChannel.Close()
	log.Println("Generating a conversation")
	fullAnalysis := ""
	summary := ""
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go s.getSentimentAnalysis(convoChannel, &fullAnalysis, wg)
	go s.getSummary(convoChannel, &summary, wg)
	wg.Wait()
	if fullAnalysis == "" || summary == "" {
		convoChannel.FullUpdate(models.ERROR, "conversationService: fullAnalysis or summary is empty")
		return
	}
	convoChannel.StatusUpdate(models.COMBINING_ANSWERS)
	temp := strings.Split(fullAnalysis, "\n")
	analysis := temp[0]
	purpose := temp[1]

	complexPrompt := s.promptGenerator.GenerateComplexPrompt(analysis, purpose, summary)
	convoChannel.StatusUpdate(models.LLM3_STARTED)
	reader, err := s.llm3.GenerateTextAsStream(complexPrompt)
	if err != nil {
		convoChannel.FullUpdate(models.ERROR, err.Error())
		return
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
}

func (s *ConversationService) getSentimentAnalysis(convoChannel ConversationChannel, result *string, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("Getting sentiment analysis")
	payload := convoChannel.Payload()
	sentimentAnalysisPrompt := s.promptGenerator.GenerateSentimentAnalysisPrompt(payload)
	convoChannel.StatusUpdate(models.LLM1_STARTED)
	var err error
	*result, err = s.llm1.GenerateText(sentimentAnalysisPrompt)
	if err != nil {
		convoChannel.FullUpdate(models.ERROR, err.Error())
		return
	}
	convoChannel.StatusUpdate(models.LLM1_FINISHED)
}

func (s *ConversationService) getSummary(convoChannel ConversationChannel, result *string, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("Getting summary")
	payload := convoChannel.Payload()
	summaryPrompt := s.promptGenerator.GenerateSummarizePrompt(payload)
	convoChannel.StatusUpdate(models.LLM2_STARTED)
	var err error
	*result, err = s.llm2.GenerateText(summaryPrompt)
	if err != nil {
		convoChannel.FullUpdate(models.ERROR, err.Error())
		return
	}
	convoChannel.StatusUpdate(models.LLM2_FINISHED)
}
