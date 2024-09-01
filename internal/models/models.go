package models

import "time"

type Message struct {
	Content    string
	Author     string
	Channel    string
	ID         string
	Platform   string
	RecievedAt time.Time
}

type MessageReaction struct {
	Message  *Message
	Reaction string
}
