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
	tuiApp := tui{
		app:   tview.NewApplication(),
		pages: tview.NewPages(),
	}

	conf := &config{
		apiClient: openLibClient,
		tui:       tuiApp,
		registry:  registerCommands(),
	}

	// TODO: uncomment once were ready to merge these things together
	// startCli(conf)

	app := conf.tui.app
	pages := conf.tui.pages
	registry := conf.registry
	commands := tview.NewList().ShowSecondaryText(false)

	for _, c := range registry {
		commands.SetTitle("Functions").SetBorder(true)
		commands.AddItem(c.name, "", 0, nil)
	}

	for _, c := range registry {
		results := conf.tui.newResultsView()
		view := c.view(conf)

		commands.SetSelectedFunc(func(i int, main string, secondary string, shortcut rune) {
			pages.SwitchToPage(c.name)
			app.SetFocus(view)
		})

		// TODO: ugh, command listItems need to be able to switch pages, and set focus to flexbox item child tview.Primitive
		// indices := commands.FindItems(c.name, "", false, true)
		// if len(indices) <= 1 {
		// 	i := indices[0]
		// 	commands.SetSelectedFunc(func(index i, c.name, "", 0) {
		//
		// 	})
		// } else {
		// 	err := fmt.Errorf("multiple commands registered with identical name")
		// 	panic(err)
		// }

		switch v := view.(type) {
		case *tview.InputField:
			v.SetDoneFunc(func(key tcell.Key) {
				if key == tcell.KeyEnter {
					searchText := v.GetText()
					// TODO: utilize the callback from the command definition
					results.SetText(fmt.Sprintf("\n\n%s", searchText))
				} else if key == tcell.KeyEsc {
					app.SetFocus(commands)
				}
			})

		}

		flexLayout := tview.NewFlex().
			AddItem(commands, 0, 1, true).
			AddItem(view, 0, 1, false).
			AddItem(results, 0, 3, false)

		pages.AddPage(c.name, flexLayout, true, true)
	}

	app.SetRoot(pages, true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func (tui *tui) newResultsView() *tview.TextView {
	results := tview.NewTextView().
		SetChangedFunc(func() {
			tui.app.Draw()
		})
	results.SetTitle("Results").SetBorder(true)

	return results
}
