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

func viewSearch(conf *config) tview.Primitive {
	search := tview.NewInputField().
		SetLabel("Title").
		SetFieldWidth(20)
	search.SetTitle("Search").SetBorder(true)

	return search
}

func resultSearch(conf *config, data []byte) tview.Primitive {
	results := tview.NewTextView().
		SetChangedFunc(func() {
			conf.tui.app.Draw()
		})
	results.SetTitle("Search Results").SetBorder(true)
	results.SetText(string(data))

	return results
}
