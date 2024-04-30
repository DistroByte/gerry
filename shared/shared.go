package shared

type Bot interface {
	Register(command string, f PluginCallFunc)

	Send(context MessageContext, message string)
}

type MessageContext struct {
	Sender  string
	Target  string
	Source  string
	Message string
}

type Config struct {
	PluginPath string
}

type PluginSetupFunc = func(bot Bot)
type PluginCallFunc = func(bot Bot, context MessageContext, arguments []string, config Config) error
type PluginRunFunc = func(bot Bot, stop <-chan struct{})
