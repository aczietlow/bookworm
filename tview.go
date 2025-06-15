package main

import (
	"github.com/rivo/tview"
)

type commandView struct {
	resultView       tview.Primitive
	updateResultView func(tview.Primitive, []byte)
	view             tview.Primitive
	updateView       func(tview.Primitive, []byte)
}

func (cv *commandView) UpdateView(data []byte) {
	cv.updateView(cv.view, data)
}

func (cv *commandView) UpdateResultView(data []byte) {
	cv.updateResultView(cv.resultView, data)
}
