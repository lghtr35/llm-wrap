package services

import (
	"github.com/google/uuid"
	"github.com/lghtr35/llm-wrap/models"
)

type ConversationChannel struct {
	channel chan models.ConversationEvent
	Event   models.ConversationEvent
}

func NewConversationChannel(payload string) ConversationChannel {
	return ConversationChannel{
		channel: make(chan models.ConversationEvent),
		Event: models.ConversationEvent{
			Conversation: models.Conversation{
				ContextID: uuid.New().String(),
				Payload:   payload,
				Response:  "",
			},
			Status: models.PREPARING_FOR_CONVERSATION,
		},
	}
}

func (c *ConversationChannel) StatusUpdate(status models.Status) {
	c.Event.Status = status
	c.channel <- c.Event
}

func (c *ConversationChannel) FullUpdate(status models.Status, response string) {
	c.Event.Status = status
	c.Event.Conversation.Response = response
	c.channel <- c.Event
}

func (c *ConversationChannel) Close() {
	close(c.channel)
}

func (c *ConversationChannel) Payload() string {
	return c.Event.Conversation.Payload
}

func (c *ConversationChannel) Channel() <-chan models.ConversationEvent {
	return c.channel
}
