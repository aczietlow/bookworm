package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/aczietlow/addToBookshelf/pkg/openlibraryapi"
)

type config struct {
	apiClient openlibraryapi.Client
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

var registry map[string]cliCommand

func startCli(conf *config) {
	scanner := bufio.NewScanner(os.Stdin)
	registry := registerCommands()
	fmt.Print("Bookworm > ")
	for scanner.Scan() {
		userInput := cleanInput(scanner.Text())

		command := userInput[0]
		if c, ok := registry[command]; ok {
			err := c.callback(conf, userInput[1:]...)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
			}
		} else {
			fmt.Print("Unknown Command\n")
		}
		fmt.Print("Bookworm > ")
	}

}

func registerCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit library api",
			callback:    commandExit,
		},
	}
}

func cleanInput(text string) []string {
	text = strings.TrimSpace(text)
	return strings.Fields(strings.ToLower(text))
}
