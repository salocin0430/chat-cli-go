package interfaces

import (
	"chat-cli/internal/domain"
	"time"
)

// MessageOutputPort is the interface that the chat service must implement
type MessageOutputPort interface {
	PublishMessage(msg *domain.Message) error
	SubscribeToChannel(channel string, handler func(*domain.Message)) error
	GetHistoricalMessages(channel string, duration time.Duration) ([]*domain.Message, error)
	Close() error
}
