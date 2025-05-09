package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/aczietlow/bookworm/pkg/openlibraryapi"
)

var registry map[string]cliCommand

type config struct {
	apiClient openlibraryapi.Client
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

func startCli(conf *config) {
	scanner := bufio.NewScanner(os.Stdin)
	registry = registerCommands()
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
		"help": {
			name:        "help",
			description: "List all available commands",
			callback:    commandHelp,
		},
		"search": {
			name:        "search",
			description: "Search open library via a solr query. search <string>",
			callback:    commandSearch,
		},
	}
}

func cleanInput(text string) []string {
	text = strings.TrimSpace(text)
	return strings.Fields(strings.ToLower(text))
}
