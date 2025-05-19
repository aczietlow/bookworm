package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/aczietlow/bookworm/pkg/openlibraryapi"
	"github.com/rivo/tview"
)

var registry map[string]cliCommand

type config struct {
	apiClient openlibraryapi.Client
	tui       tui
	registry  map[string]cliCommand
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
	view        func(*config) tview.Primitive
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
			view:        viewExit,
		},
		"help": {
			name:        "help",
			description: "List all available commands",
			callback:    commandHelp,
			view:        viewHelp,
		},
		"search": {
			name:        "search",
			description: "Search open library via a solr query. search <string>",
			callback:    commandSearch,
			view:        viewSearch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a book by providing its id",
			callback:    commandInspect,
			view:        viewInspect,
		},
	}
}

func cleanInput(text string) []string {
	text = strings.TrimSpace(text)
	return strings.Fields(strings.ToLower(text))
}
