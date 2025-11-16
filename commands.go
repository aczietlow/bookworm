package main

import "github.com/rivo/tview"

func registerCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"search": {
			name:        "search",
			description: "Search open library via a solr query. search <string>",
			callback:    commandSearch,
			view:        viewSearch,
			model:       initSearch(),
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a book by providing its id",
			callback:    commandInspect,
			view:        viewInspect,
			model:       initInspect(),
		},
		"help": {
			name:        "help",
			description: "List all available commands",
			callback:    commandHelp,
			view:        viewHelp,
		},
		"debug": {
			name:        "debug",
			description: "Useful for debuggin while building out app",
			callback:    commandDebug,
			view:        viewDebug,
		},
		"exit": {
			name:        "exit",
			description: "Exit library api",
			callback:    commandExit,
			view:        viewExit,
		},
	}
}

type model interface {
	update(*config, msg)
	view(*config) tview.Primitive
}

type msg interface{}
