package main

import (
	"fmt"

	"github.com/rivo/tview"
)

func commandHelp(conf *config, args ...string) error {
	fmt.Print("Bookworm usage:\n\n")
	for _, c := range registry {
		fmt.Printf("%s: %s\n", c.name, c.description)
	}
	fmt.Println()
	return nil
}

func viewHelp(conf *config) tview.Primitive {
	return tview.NewBox().SetTitle("help").SetBorder(true)
}
