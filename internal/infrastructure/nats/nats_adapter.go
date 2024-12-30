package nats

import (
	"chat-cli/internal/domain"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

type NatsAdapter struct {
	conn       *nats.Conn
	js         nats.JetStreamContext
	handlers   map[string]func(*domain.Message)
	streamName string
}

func NewNatsAdapter(url string) (*NatsAdapter, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	// Create JetStream context
	js, err := conn.JetStream()
	if err != nil {
		return nil, fmt.Errorf("error creating JetStream context: %v", err)
	}

	adapter := &NatsAdapter{
		conn:       conn,
		js:         js,
		handlers:   make(map[string]func(*domain.Message)), // map to store handlers for each channel
		streamName: "CHAT",
	}

	// Configure the stream
	if err := adapter.setupStream(); err != nil {
		return nil, err
	}

	return adapter, nil
}

func (a *NatsAdapter) setupStream() error {
	// Configure stream with retention of 1 hour
	streamConfig := &nats.StreamConfig{
		Name:      a.streamName,
		Subjects:  []string{"chat.*"},
		Retention: nats.LimitsPolicy,
		MaxAge:    1 * time.Hour,
		Storage:   nats.FileStorage, // store messages in files
		Replicas:  1,
		Discard:   nats.DiscardOld,
	}

	// Try to create the stream
	_, err := a.js.StreamInfo(a.streamName)
	if err != nil {
		// If it doesn't exist, create it
		_, err = a.js.AddStream(streamConfig)
	} else {
		// If it exists, update it
		_, err = a.js.UpdateStream(streamConfig)
	}

	if err != nil {
		return fmt.Errorf("error configuring stream: %v", err)
	}

	fmt.Printf("Stream %s configured correctly\n", a.streamName)
	return nil
}

func (a *NatsAdapter) PublishMessage(msg *domain.Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// Publicar usando JetStream
	_, err = a.js.Publish(fmt.Sprintf("chat.%s", msg.Channel), data)
	return err
}

func (a *NatsAdapter) SubscribeToChannel(channel string, handler func(*domain.Message)) error {
	a.handlers[channel] = handler

	// Use ephemeral subscription
	_, err := a.js.Subscribe(
		fmt.Sprintf("chat.%s", channel),
		func(m *nats.Msg) {
			var msg domain.Message
			if err := json.Unmarshal(m.Data, &msg); err != nil {
				return
			}
			handler(&msg)
		},
		nats.DeliverNew(), // deliver only new messages
		nats.AckNone(),    // we don't need confirmation
	)

	return err
}

func (a *NatsAdapter) GetHistoricalMessages(channel string, duration time.Duration) ([]*domain.Message, error) {
	var messages []*domain.Message
	msgChan := make(chan *domain.Message, 100)
	done := make(chan bool)

	// Create a temporary consumer to get historical messages
	sub, err := a.js.Subscribe(
		fmt.Sprintf("chat.%s", channel),
		func(m *nats.Msg) {
			var msg domain.Message
			if err := json.Unmarshal(m.Data, &msg); err != nil {
				return
			}
			msgChan <- &msg
		},
		nats.DeliverAll(), // delivery all messages from the past
		nats.StartTime(time.Now().Add(-duration)),
	)
	if err != nil {
		return nil, fmt.Errorf("error subscribing: %v", err)
	}
	defer sub.Unsubscribe() // unsubscribe at the end

	// Wait for messages with timeout (2 seconds) to avoid blocking the main thread
	go func() {
		time.Sleep(2 * time.Second)
		done <- true
	}()

	// Collect messages until timeout
	collecting := true
	for collecting {
		select {
		case msg := <-msgChan:
			messages = append(messages, msg)
		case <-done:
			collecting = false
		}
	}

	return messages, nil
}

func (a *NatsAdapter) Close() error {
	// First close the connection
	a.conn.Close()
	return nil
}
