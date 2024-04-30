package shared

type Bot interface {
	Register(command string, f PluginCallFunc)

	Send(target string, message string)
}

type MessageContext struct {
	Sender  string
	Target  string
	Message string
}

type Config struct {
	PluginPath string
}

type PluginSetupFunc = func(bot Bot)
type PluginCallFunc = func(bot Bot, context MessageContext, arguments []string, config Config) error
type PluginRunFunc = func(bot Bot, stop <-chan struct{})
