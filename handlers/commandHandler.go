package handlers

import (
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lghtr35/llm-wrap/models"
	"github.com/lghtr35/llm-wrap/services"
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
	log.Println("Received a command")
	var request models.CommandRequest
	err := c.ShouldBind(&request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	conversationChannel := services.NewConversationChannel(request.Prompt)

	go h.conversationService.GenerateConversation(conversationChannel)

	c.Stream(func(w io.Writer) bool {
		if event, ok := <-conversationChannel.Channel(); ok {
			log.Println("Streaming a conversation")
			c.SSEvent("status", event.Status.String())
			c.SSEvent("conversation", event.Conversation)
			return true
		}

		return false
	})
}
