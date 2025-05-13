package main

import (
	"fmt"
	"strings"
)

func commandSearch(conf *config, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("Please provide a a string to query")
	}

	results, err := conf.apiClient.SearchQuery(solrQueryString(args))
	if err != nil {
		return err
	}

	for _, book := range results {
		fmt.Printf("%s | %s by %s\n", book.Key, book.Title, book.AuthorName[0])
	}

	return nil
}

func solrQueryString(q []string) string {
	str := q[0]

	for _, s := range q[1:] {
		str += " " + s
	}

	return strings.ReplaceAll(str, " ", "+")
}
