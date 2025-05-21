package main

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

func commandSearch(conf *config, args ...string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("Please provide a a string to query")
	}

	searchText := solrQueryString(args[0])
	results, err := conf.apiClient.SearchQuery(searchText)
	if err != nil {
		return "", err
	}

	output := ""
	for _, book := range results {
		output += fmt.Sprintf("%s | %s by %s\n", book.Key, book.Title, book.AuthorName[0])

	}

	return output, nil
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
