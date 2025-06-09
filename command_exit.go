package main

import (
	"github.com/rivo/tview"
)

func commandExit(conf *config, args ...string) ([]byte, error) {
	conf.tui.app.Stop()
	return nil, nil
}

func viewExit(conf *config) tview.Primitive {
	return nil
}

func resultExit(conf *config, data []byte) tview.Primitive {
	return nil
}
