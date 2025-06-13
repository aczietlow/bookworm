package main

import (
	"github.com/rivo/tview"
)

func commandExit(conf *config, args ...string) ([]byte, error) {
	conf.tui.app.Stop()
	return nil, nil
}

func exitView(conf *config) tview.Primitive {
	return tview.NewBox().SetTitle("Exit").SetBorder(true)
}

func exitResult(conf *config) tview.Primitive {
	return nil
}

func newExitCommandView(conf *config) *commandView {
	return &commandView{
		view:             exitView(conf),
		updateView:       nil,
		resultView:       nil,
		updateResultView: nil,
	}
}
