package main

import (
	"fmt"
	"os"
)

var commands = map[string]string{
	"help":    "helps you out when in dire need of information",
	"version": "outputs the client version",
}

const (
	version = "0.0.1dev"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		default:
			fmt.Printf("unknown command %s\n", os.Args[1])
			os.Exit(1)

		case "help":
			RunHelp()
		case "version":
			RunVersion()
		}
	} else {
		RunHelp()
	}
}

func RunHelp() {
	fmt.Println("Usage: travis COMMAND ...\n\nAvailable commands:\n")
	for command, description := range commands {
		fmt.Printf("\t%s\t%s\n",
			ColorCommand(command),
			ColorInfo(description))
	}
	fmt.Printf("\nrun `%s help COMMAND` for more infos\n", os.Args[0])
}

func RunVersion() {
	fmt.Println(version)
}

func ColorCommand(s string) string {
	return fmt.Sprintf("\x1b[1m%s\x1b[0m", s)
}

func ColorInfo(s string) string {
	return fmt.Sprintf("\x1b[33m%s\x1b[0m", s)
}
