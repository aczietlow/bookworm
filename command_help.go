package main

import (
	"fmt"

	"github.com/rivo/tview"
)

func commandHelp(conf *config, args ...string) ([]byte, error) {
	output := fmt.Sprint("Bookworm usage:\n\n")
	for _, c := range conf.registry {
		output += fmt.Sprintf("%s: %s\n", c.name, c.description)
	}

	return []byte(output), nil
}

func helpView(conf *config) tview.Primitive {
	return tview.NewBox().SetTitle("help").SetBorder(true)
}

func helpResultView(conf *config) tview.Primitive {
	results := tview.NewTextView().
		SetChangedFunc(func() {
			conf.tui.app.Draw()
		})
	results.SetTitle("Bookworm Help").SetBorder(true)
	return results
}
func updateHelpResultView(t tview.Primitive, data []byte) {
	if tv, ok := t.(*tview.TextView); ok {
		tv.SetText(string(data))
	}
}

func newHelpCommandView(conf *config) *commandView {
	return &commandView{
		view:             helpView(conf),
		updateView:       nil,
		resultView:       helpResultView(conf),
		updateResultView: updateHelpResultView,
	}
}
