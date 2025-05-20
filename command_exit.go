package main

import (
	"github.com/rivo/tview"
)

func commandExit(conf *config, args ...string) error {
	conf.tui.app.Stop()
	return nil
}

func viewExit(conf *config) tview.Primitive {
	return tview.NewBox().SetTitle("Exit").SetBorder(true)
}
