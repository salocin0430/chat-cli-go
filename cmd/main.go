package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"chat-cli/internal/application"
	"chat-cli/internal/infrastructure/cli"
	"chat-cli/internal/infrastructure/nats"
)

func main() {
	natsURL := flag.String("nats", "nats://localhost:4222", "NATS server URL")
	channel := flag.String("channel", "general", "Chat channel name") // channel to join
	username := flag.String("name", "anonymous", "Your username")     // username to use
	flag.Parse()

	// Initialize NATS adapter
	natsAdapter, err := nats.NewNatsAdapter(*natsURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer natsAdapter.Close()

	// Initialize chat service
	chatService := application.NewChatService(natsAdapter, *username)

	// Initialize CLI
	cliAdapter := cli.NewCliAdapter(chatService)

	// IMPORTANT: Set the handler BEFORE joining the channel because it will be called when a new message is received
	chatService.SetMessageHandler(cliAdapter.DisplayMessage)

	// Join the channel
	if err := chatService.JoinChannel(*channel); err != nil {
		log.Fatalf("Error joining the channel: %v", err)
	}

	// Handle termination signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		natsAdapter.Close()
		os.Exit(0)
	}()

	// Start CLI
	if err := cliAdapter.Start(); err != nil {
		log.Fatalf("Error in CLI: %v", err)
	}
}
