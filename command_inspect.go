package main

import (
	"fmt"
	"time"

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
	c := conf.registry["inspect"]
	if search, ok := cv.view.(*tview.InputField); ok {
		search.SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				searchText := search.GetText()
				done := make(chan struct{})
				var data []byte
				var err error

				go func() {
					var callbackErr error
					data, callbackErr = c.callback(conf, searchText)
					if callbackErr != nil {
						err = callbackErr
					}
					close(done)
				}()

				if err != nil {
					panic(err)
				}

				// Spin to win
				go func() {
					spinner := `-/|\`
					i := 0
					for {
						select {
						case <-done:
							return
						default:
							r := spinner[i%len(spinner)]
							t.app.QueueUpdateDraw(func() {
								cv.UpdateResultView([]byte{r})
							})
							i++
							time.Sleep(100 * time.Millisecond)
						}
					}
				}()

				go func() {
					<-done
					t.app.QueueUpdateDraw(func() {
						if err != nil {
							cv.UpdateResultView([]byte("Error: " + err.Error()))
						} else {
							cv.UpdateResultView(data)
						}
					})
				}()

				// cv.UpdateResultView(data)
				t.app.SetFocus(cv.resultView)
			} else if key == tcell.KeyEsc {
				t.app.SetFocus(t.commands)
			}
		})
	}

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
