package main

import (
	"github.com/gdamore/tcell/v2"
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

// Reusable method to set navigation to match vim motions.
func setTviewInputMethod(event *tcell.EventKey) *tcell.EventKey {
	if event.Rune() == 'j' {
		return tcell.NewEventKey(tcell.KeyDown, rune(0), tcell.ModNone)
	} else if event.Rune() == 'k' {
		return tcell.NewEventKey(tcell.KeyUp, rune(0), tcell.ModNone)
	}
	return event
}
