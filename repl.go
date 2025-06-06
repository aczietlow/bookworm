package main

import (
	"fmt"
	"sort"

	"github.com/aczietlow/bookworm/pkg/openlibraryapi"
	"github.com/gdamore/tcell/v2"
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
	callback    func(*config, ...string) (string, error)
	view        func(*config) tview.Primitive
	result      func(*config) tview.Primitive
}

func startCli(conf *config) {
	startTuiApp(conf)
	if err := conf.tui.app.Run(); err != nil {
		panic(err)
	}
}

func startTuiApp(conf *config) {
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

	commands := tview.NewList().ShowSecondaryText(false)
	commands.SetTitle("Functions").SetBorder(true).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'j' {
			return tcell.NewEventKey(tcell.KeyDown, rune(0), tcell.ModNone)
		} else if event.Rune() == 'k' {
			return tcell.NewEventKey(tcell.KeyUp, rune(0), tcell.ModNone)
		}
		return event
	})

	for i, k := range registryOrder {
		c := registry[k]
		commands.AddItem(c.name, "", 0, nil)
		view := c.view(conf)
		views = append(views, &view)
		var results tview.Primitive

		// Look for custom result views
		if c.result != nil {
			r := c.result(conf)
			results = r
		} else {
			results = conf.tui.NewResultTextView()
		}

		switch v := view.(type) {
		case *tview.InputField:
			v.SetDoneFunc(func(key tcell.Key) {
				if key == tcell.KeyEnter {
					searchText := v.GetText()
					result, err := c.callback(conf, searchText)
					if err != nil {
						panic(err)
					}

					if r, ok := results.(*tview.TextView); ok {
						r.SetText(fmt.Sprintf("%s", result))
						results = r
					}

				} else if key == tcell.KeyEsc {
					app.SetFocus(commands)
				}
			})
		// This is the default primatitive, if there is no interactivity just call the command callback immediately.
		case *tview.Box:
			v.SetFocusFunc(func() {
				result, err := c.callback(conf)
				if err != nil {
					panic(err)
				}
				if r, ok := results.(*tview.TextView); ok {
					r.SetText(fmt.Sprintf("%s", result))
					results = r
				}
				app.SetFocus(commands)
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
}

func (tui *tui) NewResultTextView() *tview.TextView {
	results := tview.NewTextView().
		SetChangedFunc(func() {
			tui.app.Draw()
		})
	results.SetTitle("Results").SetBorder(true)
	return results
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
		"debug": {
			name:        "debug",
			description: "Useful for debuggin while building out app",
			callback:    commandDebug,
			view:        viewDebug,
			result:      resultDebug,
		},
	}
}
