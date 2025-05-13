package main

import "fmt"

func commandInspect(conf *config, args ...string) error {
	id := args[0]

	book, err := conf.apiClient.GetBookById(id)
	if err != nil {
		return err
	}

	fmt.Printf("Title: %s\n", book.Title)
	fmt.Printf("Subtitle: %s\n", book.Subtitle)
	fmt.Printf("Cover: %s\n", book.Cover)

	return nil
}
