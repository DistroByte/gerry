package bot

type Message interface {
	GetContent() string
	GetSenderID() string
	GetChannelID() string
	GetPlatform() string
}

type Command interface {
	Execute(bot *Bot, msg Message, args []string) error
	Help() string
	Name() string
}
