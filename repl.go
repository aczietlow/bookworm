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
	registry  map[string]*cliCommand
}

type cliCommand struct {
	name          string
	description   string
	callback      func(*config, ...string) ([]byte, error)
	getView       func(*config) tview.Primitive
	view          tview.Primitive
	getResultView func(*config, []byte) tview.Primitive
	resultView    tview.Primitive
}

func startCli(conf *config) {
	startTuiApp(conf)
	if err := conf.tui.app.Run(); err != nil {
		panic(err)
	}
}

func startTuiApp(conf *config) {
	app := conf.tui.app
	registry := conf.registry
	pages := conf.tui.pages
	registryOrder := conf.tui.registryOrder

	for k := range registry {
		registryOrder = append(registryOrder, k)
	}
	// TODO: Come up with a better sort algo.
	sort.Sort(sort.Reverse(sort.StringSlice(registryOrder)))

	commands := tview.NewList().ShowSecondaryText(false)
	commands.SetTitle("Functions").SetBorder(true).SetInputCapture(setTviewInputMethod)
	for i, k := range registryOrder {
		c := registry[k]
		commands.AddItem(c.name, "", 0, nil)
		c.view = c.getView(conf)
		c.resultView = c.getResultView(conf, nil)

		flexLayout := tview.NewFlex().
			AddItem(commands, 0, 1, true).
			AddItem(c.view, 0, 1, false).
			AddItem(c.resultView, 0, 3, false)

		// Add interaction callbacks
		switch v := c.view.(type) {
		case *tview.InputField:
			v.SetDoneFunc(func(key tcell.Key) {
				if key == tcell.KeyEnter {
					searchText := v.GetText()
					data, err := c.callback(conf, searchText)
					if err != nil {
						panic(err)
					}

					// TODO: Attach behavior to jump from results panes to view panes
					flexLayout.RemoveItem(c.resultView)
					c.resultView = nil
					c.resultView = c.getResultView(conf, data)

					// TODO: attaching the results navigation here feels like a bit of a hack
					if rv, ok := c.resultView.(*tview.TextView); ok {
						rv.SetDoneFunc(func(key tcell.Key) {
							if key == tcell.KeyEsc {
								app.SetFocus(commands)
							}
							app.SetFocus(c.view)
						})
					}
					// end hack

					flexLayout.AddItem(c.resultView, 0, 3, false)
					app.SetFocus(c.resultView)
				} else if key == tcell.KeyEsc {
					app.SetFocus(commands)
				}
			})

		// This is the default primatitive, if there is no interactivity just call the command callback immediately.
		case *tview.Box:
			v.SetFocusFunc(func() {
				data, err := c.callback(conf)
				if err != nil {
					panic(err)
				}

				switch resultType := c.resultView.(type) {
				case *tview.TextView:
					resultType.SetText(fmt.Sprintf("%s", data))
					c.resultView = resultType
					app.SetFocus(commands)
				case *tview.TreeView:
					resultType.SetDoneFunc(func(key tcell.Key) {
						if key == tcell.KeyEsc {
							app.SetFocus(commands)
						}
					})
					app.SetFocus(c.resultView)
				}
			})
		}

		pages.AddPage(c.name, flexLayout, true, i == 0)
	}

	// TODO: This isn't working....
	// I think the problem is when we write the data to the result primitive, we're actually creating a new instance of the result
	// Set Navigation Controls for Results
	for _, k := range registryOrder {
		c := registry[k]
		switch resultView := c.resultView.(type) {
		case *tview.TextView:
			resultView.SetDoneFunc(func(key tcell.Key) {
				if key == tcell.KeyEsc {
					app.SetFocus(commands)
				}
				app.SetFocus(c.view)
			})
		}
	}

	commands.SetSelectedFunc(func(i int, main string, secondary string, shortcut rune) {
		name := registryOrder[i]
		pages.SwitchToPage(name)
		app.SetFocus(registry[name].view)
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

func setTviewInputMethod(event *tcell.EventKey) *tcell.EventKey {
	if event.Rune() == 'j' {
		return tcell.NewEventKey(tcell.KeyDown, rune(0), tcell.ModNone)
	} else if event.Rune() == 'k' {
		return tcell.NewEventKey(tcell.KeyUp, rune(0), tcell.ModNone)
	}
	return event
}

func registerCommands() map[string]*cliCommand {
	return map[string]*cliCommand{
		"exit": {
			name:          "exit",
			description:   "Exit library api",
			callback:      commandExit,
			getView:       viewExit,
			getResultView: resultExit,
		},
		"help": {
			name:          "help",
			description:   "List all available commands",
			callback:      commandHelp,
			getView:       viewHelp,
			getResultView: resultHelp,
		},
		"search": {
			name:          "search",
			description:   "Search open library via a solr query. search <string>",
			callback:      commandSearch,
			getView:       viewSearch,
			getResultView: resultSearch,
		},
		"inspect": {
			name:          "inspect",
			description:   "Inspect a book by providing its id",
			callback:      commandInspect,
			getView:       viewInspect,
			getResultView: resultInspect,
		},
		"debug": {
			name:          "debug",
			description:   "Useful for debuggin while building out app",
			callback:      commandDebug,
			getView:       viewDebug,
			getResultView: resultDebug,
		},
	}
}
