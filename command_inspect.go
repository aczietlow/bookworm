package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func commandInspect(conf *config, args ...string) ([]byte, error) {
	id := args[0]

	book, err := conf.apiClient.GetBookById(id)
	if err != nil {
		return nil, err
	}

	output := ""
	// fmt.Printf("%+v\n", book)
	output += fmt.Sprintf("Title: %s\n", book.Title)
	output += fmt.Sprintf("Subtitle: %s\n", book.Subtitle)
	output += fmt.Sprintf("Authors: %+v\n", book.Authors)
	output += fmt.Sprintf("Source: %s\n", book.Source)
	output += fmt.Sprintf("Description %s\n", book.Description)
	output += fmt.Sprintf("Genre: %+v\n", book.Genre)
	output += fmt.Sprintf("Cover: %s\n", book.Cover)
	output += fmt.Sprintf("ISBN: %s\n", book.ISBN)

	return []byte(output), nil
}

func inspectView(conf *config) tview.Primitive {
	search := tview.NewInputField().
		SetLabel("Open Library ID").
		SetFieldWidth(20)
	search.SetTitle("Inspect").SetBorder(true)

	return search
}

func updateInspectView(t tview.Primitive, data []byte) {
	if tv, ok := t.(*tview.InputField); ok {
		tv.SetText(string(data))
	}
}

func inspectResultView(conf *config) tview.Primitive {
	results := tview.NewTextView()

	results.SetTitle("Search Results").SetBorder(true)
	return results
}

func updateInspectResultView(t tview.Primitive, data []byte) {
	if tv, ok := t.(*tview.TextView); ok {
		tv.SetText(string(data))
	}
}

func newInspectCommandView(conf *config) *commandView {
	cv := &commandView{
		view:             inspectView(conf),
		updateView:       updateInspectView,
		resultView:       inspectResultView(conf),
		updateResultView: updateInspectResultView,
	}

	return cv
}

func attachInspectBehaviors(cv *commandView, conf *config) {
	t := conf.tui
	if results, ok := cv.resultView.(*tview.TextView); ok {
		results.SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEsc {
				t.app.SetFocus(t.commands)
			} else if key == tcell.KeyEnter {
				t.app.SetFocus(cv.resultView)
			}
			t.app.SetFocus(cv.view)
		})
		// TODO:I'm pretty sure this isn't really needed
		// results.SetChangedFunc(func() {
		// 	conf.tui.app.Draw()
		// })
	}

}
