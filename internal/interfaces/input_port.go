package interfaces

// ChatInputPort is the interface that the chat service must implement
type ChatInputPort interface {
	SendMessage(content string) error
	JoinChannel(channel string) error
	LeaveChannel() error
}
