package models

type Message struct {
	Content  string
	Author   string
	Channel  string
	ID       string
	Platform string
}

type MessageReaction struct {
	Message  *Message
	Reaction string
}
