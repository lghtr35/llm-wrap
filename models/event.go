package models

type Status int

const (
	PREPARING_FOR_CONVERSATION Status = iota + 1
	LLM1_STARTED
	LLM1_FINISHED
	LLM2_STARTED
	LLM2_FINISHED
	COMBINING_ANSWERS
	LLM3_STARTED
	LLM3_STREAMING
	LLM3_FINISHED
	ERROR
)

func (s Status) String() string {
	switch s {
	case PREPARING_FOR_CONVERSATION:
		return "Preparing for conversation"
	case LLM1_STARTED:
		return "LLM1 started"
	case LLM1_FINISHED:
		return "LLM1 finished"
	case LLM2_STARTED:
		return "LLM2 started"
	case LLM2_FINISHED:
		return "LLM2 finished"
	case COMBINING_ANSWERS:
		return "Combining answers"
	case LLM3_STARTED:
		return "LLM3 started"
	case LLM3_STREAMING:
		return "LLM3 streaming"
	case LLM3_FINISHED:
		return "LLM3 finished"
	case ERROR:
		return "Error"
	default:
		return "Unknown status"
	}
}

type ConversationEvent struct {
	Status Status
	Conversation
}

type Conversation struct {
	ContextID string
	Payload   string
	Response  string
}
