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

func viewHelp(conf *config) tview.Primitive {
	return tview.NewBox().SetTitle("help").SetBorder(true)
}

func resultHelp(conf *config, data []byte) tview.Primitive {
	results := tview.NewTextView().
		SetChangedFunc(func() {
			conf.tui.app.Draw()
		})
	results.SetTitle("Search Results").SetBorder(true)
	results.SetText(string(data))
	return results
}
