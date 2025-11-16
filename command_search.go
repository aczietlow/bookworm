package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/rivo/tview"
)

func commandSearch(conf *config, args ...string) (string, error) {
	// Commands should really be self contained.
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
		authorName := "N/A"
		if len(book.AuthorName) > 0 {
			authorName = book.AuthorName[0]
		}
		output += fmt.Sprintf("%s | %s by %s\n", book.Key, book.Title, authorName)
	}

	return output, nil
}

func solrQueryString(q string) string {
	return strings.ReplaceAll(q, " ", "+")
}

func viewSearch(conf *config) tview.Primitive {
	search := tview.NewInputField().
		SetLabel("Titlexx").
		SetFieldWidth(20)
	search.SetTitle("Search").SetBorder(true)

	return search
}

type searchResult struct {
	id     string
	title  string
	author string
}

type searchModel struct {
	query       string
	results     []searchResult
	tview       tview.Primitive
	tviewResult tview.Primitive
}

type searchQuery string

func initSearch() model {
	search := tview.NewInputField().
		SetLabel("Book Title").
		SetFieldWidth(30)
	search.SetTitle("Search").SetBorder(true)

	return searchModel{
		tview: search,
	}
}

func (sm searchModel) update(conf *config, message msg) {
	switch msg := message.(type) {
	case searchQuery:
		// TODO: why did I create a custom type here?
		queryText := string(msg)
		sm.query = queryText

		searchText := solrQueryString(queryText)
		results, err := conf.apiClient.SearchQuery(searchText)
		if err != nil {
			log.Println("err:\n", err)
		}

		for _, book := range results {
			authorName := "N/A"
			if len(book.AuthorName) > 0 {
				authorName = book.AuthorName[0]
			}
			sm.results = append(sm.results, searchResult{book.Key, book.Title, authorName})
		}
	}
}

func (sm searchModel) view(conf *config) tview.Primitive {
	return sm.tview
}
