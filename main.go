package main

import (
	"fmt"
	"sort"
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

	// TODO: This might be the dumbest way of tracking this state.
	var views []*tview.Primitive
	var registryOrder []string
	for k := range registry {
		registryOrder = append(registryOrder, k)
	}

	// TODO: Come up with a better sort algo.
	sort.Sort(sort.Reverse(sort.StringSlice(registryOrder)))
	// slices.Sort(registryOrder)

	commands := tview.NewList().ShowSecondaryText(false)
	commands.SetTitle("Functions").SetBorder(true)

	for i, k := range registryOrder {
		c := registry[k]
		commands.AddItem(c.name, "", 0, nil)
		results := conf.tui.newResultsView()
		view := c.view(conf)
		views = append(views, &view)

		switch v := view.(type) {
		case *tview.InputField:
			v.SetDoneFunc(func(key tcell.Key) {
				if key == tcell.KeyEnter {
					searchText := v.GetText()
					// TODO: utilize the callback from the command definition
					results.SetText(fmt.Sprintf("%s", searchText))
				} else if key == tcell.KeyEsc {
					app.SetFocus(commands)
				}
			})

		}

		flexLayout := tview.NewFlex().
			AddItem(commands, 0, 1, true).
			AddItem(view, 0, 1, false).
			AddItem(results, 0, 3, false)

		pages.AddPage(c.name, flexLayout, true, i == 0)
	}

	commands.SetSelectedFunc(func(i int, main string, secondary string, shortcut rune) {
		name := registryOrder[i]
		pages.SwitchToPage(name)

		view := views[i]
		app.SetFocus(*view)
	})

	app.SetRoot(pages, true).SetFocus(pages)
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
