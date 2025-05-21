package main

import (
	"fmt"

	"github.com/rivo/tview"
)

func commandInspect(conf *config, args ...string) (string, error) {
	id := args[0]

	book, err := conf.apiClient.GetBookById(id)
	if err != nil {
		return "", err
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

	return output, nil
}

func viewInspect(conf *config) tview.Primitive {
	search := tview.NewInputField().
		SetLabel("Open Library ID").
		SetFieldWidth(20)
	search.SetTitle("Inspect").SetBorder(true)

	return search
}
