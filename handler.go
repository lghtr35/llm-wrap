package main

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lghtr35/llm-wrapping/models"
	"github.com/lghtr35/llm-wrapping/services"
)

type CommandHandler struct {
	conversationService *services.ConversationService
}

func NewCommandHandler(vendorConfigs map[string]models.VendorConfig) *CommandHandler {
	return &CommandHandler{
		conversationService: services.NewConversationService(vendorConfigs["OpenAI"]),
	}
}

func (h *CommandHandler) Handle(c *gin.Context) {
	var request models.CommandRequest
	err := c.ShouldBind(&request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	conversationChannel := services.NewConversationChannel(request.Prompt)

	err = h.conversationService.GenerateConversation(conversationChannel)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Stream(func(w io.Writer) bool {
		if event, ok := <-conversationChannel.Channel(); ok {
			c.SSEvent("status", event.Status)
			c.SSEvent("conversation", event.Conversation)
			return true
		}

		return false
	})
}
