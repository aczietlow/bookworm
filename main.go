package main

import (
	"fmt"
	"time"

	"github.com/aczietlow/bookworm/pkg/openlibraryapi"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	// "github.com/aczietlow/bookworm/pkg/openlibraryapi"
)

// var commands = []string{
// 	"search",
// 	"help",
// 	"Quit",
// 	"Inspect",
// }

var regv2 map[string]command

// TOOD: extend existing function registry with this
type command struct {
	name        string
	description string
}

type tui struct {
	app   *tview.Application
	pages *tview.Pages
	focus tview.Primitive
}

const (
	tuiStartPage = "Search"
)

func main() {
	openLibClient := openlibraryapi.NewClient(time.Second*10, time.Minute*5)

	conf := &config{
		apiClient: openLibClient,
		tui: tui{
			app:   tview.NewApplication(),
			pages: tview.NewPages(),
		},
	}

	// TODO: uncomment once were ready to merge these things together
	// startCli(conf)

	app := conf.tui.app

	// stubbing out registry
	registry := makeCommands()

	search := tview.NewInputField().
		SetLabel("Title").
		SetFieldWidth(20)
	search.SetTitle("Search").SetBorder(true)

	// List commands
	commands := tview.NewList().ShowSecondaryText(false)
	for _, c := range registry {
		// search = c.view
		commands.AddItem(c.name, "", 0, func() {
			app.SetFocus(search)
		})
	}
	commands.AddItem("Quit", "", 'q', func() {
		app.Stop()
	})
	commands.SetTitle("Functions").SetBorder(true)

	results := tview.NewTextView().
		SetChangedFunc(func() {
			app.Draw()
		})
	results.SetTitle("Results").SetBorder(true)

	search.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			searchText := search.GetText()
			// results.Clear()
			results.SetText(fmt.Sprintf("\n\n%s", searchText))
		} else if key == tcell.KeyEsc {
			app.SetFocus(commands)
		}
	})

	flexLayout := tview.NewFlex().
		AddItem(commands, 0, 1, true).
		AddItem(search, 0, 1, false).
		AddItem(results, 0, 3, false)

	pages := tview.NewPages().AddPage(tuiStartPage, flexLayout, true, true)

	conf.tui.pages = pages
	app.SetRoot(pages, true)

	if err := app.Run(); err != nil {
		panic(err)
	}
}

// TODO: Attempting to build view for all registry commands
func makeCommands() map[string]command {
	return map[string]command{
		"Search": {
			name:        "Search",
			description: "Search Function",
		},
	}
}
