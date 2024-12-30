package application

import (
	"chat-cli/internal/domain"
	"chat-cli/internal/interfaces"
	"fmt"
	"sort"
	"time"
)

type ChatService struct {
	messageOutput  interfaces.MessageOutputPort
	username       string
	channel        string
	messageHandler func(*domain.Message)
}

func NewChatService(messageOutput interfaces.MessageOutputPort, username string) *ChatService {
	return &ChatService{
		messageOutput: messageOutput,
		username:      username,
	}
}

func (s *ChatService) SendMessage(content string) error {
	msg := &domain.Message{
		Content:   content,
		Sender:    s.username,
		Channel:   s.channel,
		Timestamp: time.Now(),
	}
	return s.messageOutput.PublishMessage(msg)
}

func (s *ChatService) JoinChannel(channel string) error {
	s.channel = channel

	fmt.Printf("\n[%s] Joining channel %s...\n", s.username, channel) // print message in the format [HH:MM:SS] Sender: Message

	// Get historical messages first
	messages, err := s.messageOutput.GetHistoricalMessages(channel, time.Hour)
	if err != nil {
		return fmt.Errorf("error getting historical messages: %v", err)
	}

	if len(messages) > 0 {
		fmt.Printf("\n=== Messages from the last hour ===\n")
		// Sort messages by timestamp
		sort.Slice(messages, func(i, j int) bool {
			return messages[i].Timestamp.Before(messages[j].Timestamp)
		})

		// Show historical messages
		for _, msg := range messages {
			if s.messageHandler != nil {
				s.messageHandler(msg)
			}
		}
		fmt.Printf("=== End of historical messages ===\n\n")
	}

	// Subscribe to new messages
	if err := s.messageOutput.SubscribeToChannel(channel, s.handleMessage); err != nil {
		return fmt.Errorf("error subscribing to new messages: %v", err)
	}

	return nil
}

func (s *ChatService) SetMessageHandler(handler func(*domain.Message)) {
	s.messageHandler = handler
}

func (s *ChatService) handleMessage(msg *domain.Message) {
	if s.messageHandler != nil {
		s.messageHandler(msg)
	}
}

func (s *ChatService) LeaveChannel() error {
	if s.channel == "" {
		return nil // No channel to leave
	}

	// exit message // TODO: add a message to the user that he has left the channel in other seminary
	msg := &domain.Message{
		Content:   "has left the channel",
		Sender:    s.username,
		Channel:   s.channel,
		Timestamp: time.Now(),
	}

	if err := s.messageOutput.PublishMessage(msg); err != nil {
		return err
	}

	s.channel = ""
	return nil
}
