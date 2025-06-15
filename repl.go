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
	name           string
	description    string
	callback       func(*config, ...string) ([]byte, error)
	getCommandView func(*config) *commandView
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

	tui.hackInit()

	tui.startTuiApp()

	return tui
}

func (t *tui) Run() {
	if err := t.app.Run(); err != nil {
		panic(err)
	}
}

func (t *tui) hackInit() {
	attachInspectBehaviors(t.commandsView["inspect"], t.appState)
	attachSearchViewBehaviors(t.commandsView["search"], t.appState)
}

func (t *tui) initTui() {
	t.buildCommandsList()
	t.buildCommandsViews()
}

func (t *tui) buildCommandsList() {
	for c := range t.appState.registry {
		t.registryOrder = append(t.registryOrder, c)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(t.registryOrder)))

	t.commands = tview.NewList().ShowSecondaryText(false)

	for _, k := range t.registryOrder {
		c := t.appState.registry[k]
		t.commands.AddItem(c.name, "", 0, nil)
	}

	t.commands.SetTitle("Functions").SetBorder(true).SetInputCapture(setTviewInputMethod)
}

func (t *tui) buildCommandsViews() {
	for _, k := range t.registryOrder {
		c := t.appState.registry[k]
		t.commandsView[c.name] = c.getCommandView(t.appState)
	}

	t.commands.SetSelectedFunc(func(i int, main string, secondary string, shortcut rune) {
		name := t.registryOrder[i]
		t.pages.SwitchToPage(name)
		t.app.SetFocus(t.commandsView[name].view)
	})

}

func (t *tui) startTuiApp() {
	conf := t.appState

	for i, k := range t.registryOrder {
		c := t.appState.registry[k]
		cv := t.commandsView[c.name]

		flexLayout := tview.NewFlex().
			AddItem(t.commands, 0, 1, true).
			AddItem(cv.view, 0, 1, false).
			AddItem(cv.resultView, 0, 3, false)

		// Add interaction callbacks
		switch v := cv.view.(type) {
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
			name:           "exit",
			description:    "Exit library api",
			callback:       commandExit,
			getCommandView: newExitCommandView,
		},
		"help": {
			name:           "help",
			description:    "List all available commands",
			getCommandView: newHelpCommandView,
			callback:       commandHelp,
		},
		"search": {
			name:           "search",
			description:    "Search open library via a solr query. search <string>",
			callback:       commandSearch,
			getCommandView: newSearchCommandView,
		},
		"inspect": {
			name:           "inspect",
			description:    "Inspect a book by providing its id",
			callback:       commandInspect,
			getCommandView: newInspectCommandView,
		},
		"debug": {
			name:           "debug",
			description:    "Useful for debuggin while building out app",
			callback:       commandDebug,
			getCommandView: newDebugCommandView,
		},
	}
}
