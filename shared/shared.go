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

type PluginSetupFunc = func(bot Bot)
type PluginCallFunc = func(bot Bot, context MessageContext, arguments []string) error
type PluginRunFunc = func(bot Bot, stop <-chan struct{})
