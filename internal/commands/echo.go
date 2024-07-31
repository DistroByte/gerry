package commands

import "strings"

func EchoCommand(args []string) string {
	return strings.Join(args[1:], " ")
}
