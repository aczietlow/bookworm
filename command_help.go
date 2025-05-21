package main

import (
	"fmt"

	"github.com/rivo/tview"
)

func commandHelp(conf *config, args ...string) (string, error) {
	output := fmt.Sprint("Bookworm usage:\n\n")
	for _, c := range registry {
		output += fmt.Sprintf("%s: %s\n", c.name, c.description)
	}
	fmt.Println()
	return output, nil
}

func viewHelp(conf *config) tview.Primitive {
	return tview.NewBox().SetTitle("help").SetBorder(true)
}
