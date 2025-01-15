package services

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	fullAnalysis := models.SentimentAnalysisResponse{}
	summary := models.SummaryAnalysisResponse{}
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go s.getSentimentAnalysis(convoChannel, &fullAnalysis, wg)
	go s.getSummary(convoChannel, &summary, wg)
	wg.Wait()
	if fullAnalysis.Purpose == "" || fullAnalysis.Sentiment == "" || summary.Summary == "" {
		convoChannel.FullUpdate(models.ERROR, "conversationService: fullAnalysis or summary is empty")
		return
	}
	fmt.Printf("Full Analysis: %+v\nSummary: %+v\n", fullAnalysis, summary)
	convoChannel.StatusUpdate(models.COMBINING_ANSWERS)
	analysis := fullAnalysis.Sentiment
	purpose := fullAnalysis.Purpose

	complexPrompt := s.promptGenerator.GenerateComplexPrompt(analysis, purpose, summary.Summary)
	convoChannel.StatusUpdate(models.LLM3_STARTED)
	body, err := s.llm3.GenerateTextAsStream(complexPrompt, nil)
	if err != nil {
		convoChannel.FullUpdate(models.ERROR, err.Error())
		return
	}
	defer body.Close()

	messages := make(chan string)
	errChannel := make(chan error)
	go s.processLines(body, messages, errChannel)
	var result string
	for {
		select {
		case err := <-errChannel:
			if err == io.EOF {
				convoChannel.FullUpdate(models.LLM3_FINISHED, result)
				return
			}
			convoChannel.FullUpdate(models.ERROR, err.Error())
			return
		default:
		}
		select {
		case msg := <-messages:
			result += msg
			convoChannel.FullUpdate(models.LLM3_STREAMING, result)
			fmt.Println("received message", msg)
		default:
		}
	}
}

func (s *ConversationService) getSentimentAnalysis(convoChannel ConversationChannel, result *models.SentimentAnalysisResponse, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("Getting sentiment analysis")
	payload := convoChannel.Payload()
	sentimentAnalysisPrompt := s.promptGenerator.GenerateSentimentAnalysisPrompt(payload)
	convoChannel.StatusUpdate(models.LLM1_STARTED)
	var err error
	responseFormat, err := models.NewResponseFormat(result.GetSchema())
	if err != nil {
		convoChannel.FullUpdate(models.ERROR, err.Error())
		return
	}
	str, err := s.llm1.GenerateText(sentimentAnalysisPrompt, &responseFormat)
	if err != nil {
		convoChannel.FullUpdate(models.ERROR, err.Error())
		return
	}
	err = json.Unmarshal([]byte(str), result)
	if err != nil {
		convoChannel.FullUpdate(models.ERROR, err.Error())
		return
	}
	convoChannel.StatusUpdate(models.LLM1_FINISHED)
}

func (s *ConversationService) getSummary(convoChannel ConversationChannel, result *models.SummaryAnalysisResponse, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("Getting summary")
	payload := convoChannel.Payload()
	summaryPrompt := s.promptGenerator.GenerateSummarizePrompt(payload)
	convoChannel.StatusUpdate(models.LLM2_STARTED)
	var err error
	responseFormat, err := models.NewResponseFormat(result.GetSchema())
	if err != nil {
		convoChannel.FullUpdate(models.ERROR, err.Error())
		return
	}
	str, err := s.llm2.GenerateText(summaryPrompt, &responseFormat)
	if err != nil {
		convoChannel.FullUpdate(models.ERROR, err.Error())
		return
	}
	err = json.Unmarshal([]byte(str), result)
	if err != nil {
		convoChannel.FullUpdate(models.ERROR, err.Error())
		return
	}
	convoChannel.StatusUpdate(models.LLM2_FINISHED)
}

func (s *ConversationService) processLines(body io.ReadCloser, channel chan string, errChannel chan error) {
	hasErrorPrefix := false
	dataPrefix := []byte("data: ")
	errorPrefix := []byte("data: {\"error\":")
	errorBytes := make([]byte, 0)
	reader := bufio.NewReader(body)
	for {
		var message models.OpenAIStreamingResponse
		raw, err := reader.ReadBytes('\n')
		if err != nil {
			errChannel <- err
			break
		}
		if hasErrorPrefix {
			var apiError models.OpenAIErrorResponse
			err = json.Unmarshal(errorBytes, &apiError)
			if err != nil {
				continue
			}
			errChannel <- fmt.Errorf("OpenAiClient: %s", apiError.Error.Message)
			break
		}
		trimmed := bytes.TrimSpace(raw)
		if bytes.HasPrefix(trimmed, errorPrefix) {
			hasErrorPrefix = true
		}
		if !bytes.HasPrefix(trimmed, dataPrefix) || hasErrorPrefix {
			if hasErrorPrefix {
				trimmed = bytes.TrimPrefix(trimmed, dataPrefix)
			}
			errorBytes = append(errorBytes, trimmed...)

			continue
		}

		rawData := bytes.TrimPrefix(trimmed, dataPrefix)
		if string(rawData) == "[DONE]" {
			errChannel <- io.EOF
			break
		}

		err = json.Unmarshal(rawData, &message)
		if err != nil {
			errChannel <- err
		}
		channel <- message.Choices[0].Delta.Content
	}
}
