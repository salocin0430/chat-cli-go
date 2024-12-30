package cli

import (
	"bufio"
	"chat-cli/internal/domain"
	"chat-cli/internal/interfaces"
	"fmt"
	"os"
)

type CliAdapter struct {
	chatService interfaces.ChatInputPort
	scanner     *bufio.Scanner // scanner to read input from the user
}

func NewCliAdapter(chatService interfaces.ChatInputPort) *CliAdapter {
	return &CliAdapter{
		chatService: chatService,
		scanner:     bufio.NewScanner(os.Stdin), // initialize scanner to read input from the user
	}
}

func (a *CliAdapter) Start() error {
	fmt.Println("Chat started. Write your messages (Ctrl+C to exit):")

	for a.scanner.Scan() {
		message := a.scanner.Text()
		if err := a.chatService.SendMessage(message); err != nil {
			fmt.Printf("Error sending message: %v\n", err)
		}
	}

	return a.scanner.Err()
}

func (a *CliAdapter) DisplayMessage(msg *domain.Message) {
	timeStr := msg.Timestamp.Format("15:04:05")                   // format time to HH:MM:SS
	fmt.Printf("[%s] %s: %s\n", timeStr, msg.Sender, msg.Content) // print message in the format [HH:MM:SS] Sender: Message
}
