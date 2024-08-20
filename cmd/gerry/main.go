package main

import (
	"fmt"
	"os"
	"time"

	"github.com/DistroByte/gerry/internal/bot"
	"github.com/DistroByte/gerry/internal/config"
)

func main() {
	providedArgs := os.Args[1:]

	if len(providedArgs) == 0 {
		fmt.Println("no arguments provided")
		os.Exit(1)
	}

	switch providedArgs[0] {

	case "confgen":
		config.Generate(providedArgs[1])

	case "start":
		if providedArgs[1] == "" {
			fmt.Println("no config file provided")
			os.Exit(1)
		}

		config.StartTime = time.Now()
		bot.Start(providedArgs[1])
	}
}
