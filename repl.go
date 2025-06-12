package main

import (
	"fmt"
	"sort"

	"github.com/aczietlow/bookworm/pkg/openlibraryapi"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type config struct {
	apiClient openlibraryapi.Client
	tui       *tui
	registry  map[string]*cliCommand
}

type cliCommand struct {
	name          string
	description   string
	callback      func(*config, ...string) ([]byte, error)
	getView       func(*config) tview.Primitive
	getResultView func(*config, []byte) tview.Primitive
}

type commandView struct {
	resultView tview.Primitive
	view       tview.Primitive
	updateView func([]byte)
}

type tuiState struct {
	currentBook string
}

type tui struct {
	appState      *config
	app           *tview.Application
	pages         *tview.Pages
	commands      *tview.List
	registryOrder []string
	commandsView  map[string]*commandView
	tuiState      *tuiState
}

func NewTui(conf *config) *tui {
	tui := &tui{
		app:          tview.NewApplication(),
		pages:        tview.NewPages(),
		appState:     conf,
		tuiState:     &tuiState{},
		commandsView: make(map[string]*commandView),
	}
	tui.initTui()

	// TODO: refactor to remove this tight coupling.
	// Actually use DI if that's the route we want to go
	conf.tui = tui

	tui.startTuiApp()

	return tui
}

func (t *tui) Run() {
	if err := t.app.Run(); err != nil {
		panic(err)
	}
}

func (t *tui) initTui() {
	t.buildCommandsList()
}

func (t *tui) buildCommandsList() {
	for c := range t.appState.registry {
		t.registryOrder = append(t.registryOrder, c)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(t.registryOrder)))

	t.commands = tview.NewList().ShowSecondaryText(false).SetSelectedFunc(func(i int, main string, secondary string, shortcut rune) {
		name := t.registryOrder[i]
		t.pages.SwitchToPage(name)
		t.app.SetFocus(t.commandsView[name].view)
	})

	t.commands.SetTitle("Functions").SetBorder(true).SetInputCapture(setTviewInputMethod)
}

func (t *tui) startTuiApp() {
	conf := t.appState

	for i, k := range t.registryOrder {
		c := t.appState.registry[k]
		t.commands.AddItem(c.name, "", 0, nil)
		t.commandsView[c.name] = &commandView{
			view:       c.getView(conf),
			resultView: c.getResultView(conf, nil),
		}
		cv := t.commandsView[c.name]

		flexLayout := tview.NewFlex().
			AddItem(t.commands, 0, 1, true).
			AddItem(cv.view, 0, 1, false).
			AddItem(cv.resultView, 0, 3, false)

		// Add interaction callbacks
		switch v := cv.view.(type) {
		case *tview.InputField:
			v.SetDoneFunc(func(key tcell.Key) {
				if key == tcell.KeyEnter {
					searchText := v.GetText()
					data, err := c.callback(conf, searchText)
					if err != nil {
						panic(err)
					}

					flexLayout.RemoveItem(cv.resultView)
					cv.resultView = nil
					cv.resultView = c.getResultView(conf, data)

					// TODO: attaching the results navigation here feels like a bit of a hack
					if rv, ok := cv.resultView.(*tview.TextView); ok {
						rv.SetDoneFunc(func(key tcell.Key) {
							if key == tcell.KeyEsc {
								t.app.SetFocus(t.commands)
							}
							t.app.SetFocus(cv.view)
						})
					}
					if rv, ok := cv.resultView.(*tview.List); ok {
						rv.SetDoneFunc(func() {
							if key == tcell.KeyEsc {
								t.app.SetFocus(cv.view)
							}
						}).SetSelectedFunc(func(i int, main string, secondary string, shortcut rune) {
							t.tuiState.currentBook = main
							if iv, ok := t.commandsView["inspect"].view.(*tview.InputField); ok {
								iv.SetText(main)
							}
							t.app.SetFocus(t.commands)
						})
					}
					// end hack

					flexLayout.AddItem(cv.resultView, 0, 3, false)
					t.app.SetFocus(cv.resultView)
				} else if key == tcell.KeyEsc {
					t.app.SetFocus(t.commands)
				}
			})

		// This is the default primatitive, if there is no interactivity just call the command callback immediately.
		case *tview.Box:
			v.SetFocusFunc(func() {
				data, err := c.callback(conf)
				if err != nil {
					panic(err)
				}

				switch resultType := cv.resultView.(type) {
				case *tview.TextView:
					resultType.SetText(fmt.Sprintf("%s", data))
					cv.resultView = resultType
					t.app.SetFocus(t.commands)
				case *tview.TreeView:
					resultType.SetDoneFunc(func(key tcell.Key) {
						if key == tcell.KeyEsc {
							t.app.SetFocus(t.commands)
						}
					})
					t.app.SetFocus(cv.resultView)
				}
			})
		}

		t.pages.AddPage(c.name, flexLayout, true, i == 0)
	}

	t.app.SetRoot(t.pages, true).SetFocus(t.pages)
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
