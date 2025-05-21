package main

import (
	"time"

	"github.com/aczietlow/bookworm/pkg/openlibraryapi"
	"github.com/rivo/tview"
)

type tui struct {
	app   *tview.Application
	pages *tview.Pages
	focus tview.Primitive
}

func main() {
	openLibClient := openlibraryapi.NewClient(time.Second*10, time.Minute*5)
	conf := &config{
		apiClient: openLibClient,
		tui: tui{
			app:   tview.NewApplication(),
			pages: tview.NewPages(),
		},
		registry: registerCommands(),
	}

	startCli(conf)
}
