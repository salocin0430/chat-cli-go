package domain

import "time"

type Message struct {
	Content   string    `json:"content"`   // content of the message
	Sender    string    `json:"sender"`    // sender of the message
	Channel   string    `json:"channel"`   // channel of the message
	Timestamp time.Time `json:"timestamp"` // timestamp of the message
}
