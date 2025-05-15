package main

import "fmt"

func commandInspect(conf *config, args ...string) error {
	id := args[0]

	book, err := conf.apiClient.GetBookById(id)
	if err != nil {
		return err
	}

	// fmt.Printf("%+v\n", book)
	fmt.Printf("Title: %s\n", book.Title)
	fmt.Printf("Subtitle: %s\n", book.Subtitle)
	fmt.Printf("Authors: %+v\n", book.Authors)
	fmt.Printf("Source: %s\n", book.Source)
	fmt.Printf("Description %s\n", book.Description)
	fmt.Printf("Genre: %+v\n", book.Genre)
	fmt.Printf("Cover: %s\n", book.Cover)
	fmt.Printf("ISBN: %s\n", book.ISBN)

	return nil
}
