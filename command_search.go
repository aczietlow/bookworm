package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func commandSearch(conf *config, args ...string) ([]byte, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("Please provide a a string to query")
	}

	searchText := solrQueryString(args[0])
	results, err := conf.apiClient.SearchQuery(searchText)
	if err != nil {
		return nil, err
	}

	output := ""
	for _, book := range results {
		authorName := "N/A"
		if len(book.AuthorName) > 0 {
			authorName = book.AuthorName[0]
		}
		output += fmt.Sprintf("%s | %s by %s\n", extractWorkID(book.Key), book.Title, authorName)
	}

	return []byte(output), nil
}

func extractWorkID(path string) string {
	parts := strings.Split(path, "/")
	return strings.TrimSpace(parts[len(parts)-1])
}

func solrQueryString(q string) string {
	return strings.ReplaceAll(q, " ", "+")
}

func searchView(conf *config) tview.Primitive {
	search := tview.NewInputField().
		SetLabel("Title").
		SetFieldWidth(20)
	search.SetTitle("Search").SetBorder(true)

	return search
}

func updateSearchView(t tview.Primitive, data []byte) {
	if tv, ok := t.(*tview.InputField); ok {
		tv.SetText(string(data))
	}
}

func searchResultView(conf *config) tview.Primitive {
	list := tview.NewList()
	list.SetTitle("Search Results").SetBorder(true).SetInputCapture(setTviewInputMethod)

	return list
}

func updateSearchResultView(t tview.Primitive, data []byte) {
	if tv, ok := t.(*tview.List); ok {
		tv.Clear()
		results := strings.Split(string(data), "\n")
		for _, r := range results {
			text := strings.Split(r, "|")
			if len(text) > 1 {
				tv.AddItem(strings.TrimSpace(text[0]), text[1], 0, nil)
			} else if len(text) == 1 {
				tv.AddItem(strings.TrimSpace(text[0]), "", 0, nil)
			}
		}
	}
}

func newSearchCommandView(conf *config) *commandView {
	cv := &commandView{
		view:             searchView(conf),
		updateView:       updateSearchView,
		resultView:       searchResultView(conf),
		updateResultView: updateSearchResultView,
	}

	return cv
}

func spinner(d time.Duration) {
}

// TODO: Update this to use an event loop
func attachSearchViewBehaviors(cv *commandView, conf *config) {
	c := conf.registry["search"]
	t := conf.tui
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
							// Queue Update Draw lets go routines outside the main thread make changes a view
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

				t.app.SetFocus(cv.resultView)
			} else if key == tcell.KeyEsc {
				t.app.SetFocus(t.commands)
			}
		})
	}

	if list, ok := cv.resultView.(*tview.List); ok {
		list.SetDoneFunc(func() {
			t.app.SetFocus(cv.view)
		}).SetSelectedFunc(func(i int, main string, secondary string, shortcut rune) {
			// t.tuiState.currentBook = main
			// TODO: update state & build message to let the inspect command know it needs to update
			if iv, ok := t.commandsView["inspect"].view.(*tview.InputField); ok {
				iv.SetText(main)
			}
			t.app.SetFocus(t.commands)
		})

	}

}
