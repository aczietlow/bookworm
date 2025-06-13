package main

import (
	"fmt"
	"strings"

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
		output += fmt.Sprintf("%s | %s by %s\n", book.Key, book.Title, authorName)
	}

	return []byte(output), nil
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
		results := strings.Split(string(data), "\n")
		for _, r := range results {
			text := strings.Split(r, "|")
			if len(text) > 1 {
				tv.AddItem(text[0], text[1], 0, nil)
			}
		}
	}
}

func newSearchCommandView(conf *config) *commandView {
	return &commandView{
		view:             searchView(conf),
		updateView:       updateSearchView,
		resultView:       searchResultView(conf),
		updateResultView: updateSearchResultView,
	}
}
